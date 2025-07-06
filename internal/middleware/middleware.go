package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"urlshort/internal/auth"
	"urlshort/internal/models"
)

func RedirectIfAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err == nil && cookie != nil && cookie.Value != "" {
			return c.Redirect(http.StatusSeeOther, "/profile")
		}

		return next(c)
	}
}

func AddUsername(secret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				c.Set("username", models.UserToken{})
				return next(c)
			}

			tokenValue := cookie.Value

			claims, ok := auth.ValidateJWT(tokenValue, secret)
			if !ok {
				c.Set("username", models.UserToken{})
				return next(c)
			}

			username, ok := claims["username"]
			if !ok {
				c.Set("username", models.UserToken{})
				return next(c)
			}

			c.Set("username", models.UserToken{
				Username: username.(string),
				Ok:       true,
			})

			return next(c)
		}
	}
}
