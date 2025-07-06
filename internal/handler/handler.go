package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"urlshort/internal/interfaces"
	"urlshort/internal/models"
)

type Handler struct {
	Database       interfaces.Database
	Cache          interfaces.Cache
	IsCacheEnabled bool
	JWTSecret      []byte
}

func NewHandlers(database interfaces.Database, cache interfaces.Cache, isCacheEnabled bool, JWTSecret []byte) *Handler {
	return &Handler{
		Database:       database,
		Cache:          cache,
		IsCacheEnabled: isCacheEnabled,
		JWTSecret:      JWTSecret,
	}
}

func (h *Handler) MainPage(c echo.Context) error {
	var IsAuthenticated bool
	user := c.Get("username").(models.UserToken)
	if user.Ok {
		IsAuthenticated = true
	}

	data := map[string]any{
		"IsAuthenticated": IsAuthenticated,
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (h *Handler) Profile(c echo.Context) error {
	var username string
	user := c.Get("username").(models.UserToken)
	if user.Ok {
		username = user.Username
	}

	links, err := h.Database.GetLinksByUser(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	}

	data := map[string]any{
		"Username": username,
		"Links":    links,
		"BaseURL":  c.Scheme() + "://" + c.Request().Host,
	}

	return c.Render(http.StatusOK, "profile.html", data)
}
