package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"urlshort/internal/config"
	"urlshort/internal/database"
	"urlshort/internal/handler"
	cmiddleware "urlshort/internal/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// Setup config, db, cache
	cfg := config.LoadConfig()
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	var cache *database.Cache
	if cfg.Redis.Enabled {
		cache, err = database.NewRedis(cfg.Redis)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Setup middlewares, logging
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	e := echo.New()
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Output: file,
	// }))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Template rendering
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	// Validator
	e.Validator = &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	// Handlers
	handlers := handler.NewHandlers(db, cache, cfg.Redis.Enabled, cfg.JwtSecret)
	AddUsernameMiddleware := cmiddleware.AddUsername(cfg.JwtSecret)

	// Auth handlers
	e.File("/register", "static1234/register.html", cmiddleware.RedirectIfAuthenticated)
	e.File("/login", "static1234/login.html", cmiddleware.RedirectIfAuthenticated)
	e.POST("/register", handlers.Register, cmiddleware.RedirectIfAuthenticated)
	e.POST("/login", handlers.Login, cmiddleware.RedirectIfAuthenticated)
	e.GET("/logout", handlers.Logout)

	// Profile and API handlers
	e.GET("/profile", handlers.Profile, AddUsernameMiddleware)
	api := e.Group("/api")
	api.POST("/short", handlers.ShortUrl, AddUsernameMiddleware)
	api.DELETE("/delete/:alias", handlers.DeleteLink, AddUsernameMiddleware)

	// Default handlers
	e.GET("/:alias", handlers.RedirectHandler)
	e.GET("/", handlers.MainPage, AddUsernameMiddleware)
	e.GET("/favicon.ico", func(c echo.Context) error { return c.NoContent(200) })

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGTERM, // systemctl stop
	)
	defer stop()

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()

	db.CloseDatabase()
	if cfg.Redis.Enabled {
		cache.CloseCache()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	file.Close()
}
