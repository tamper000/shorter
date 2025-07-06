package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"urlshort/internal/auth"
	"urlshort/internal/models"
	"urlshort/internal/utils"
)

func (h *Handler) Register(c echo.Context) error {
	info := new(models.User)

	if err := c.Bind(info); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrSetLoginAndPassword)
	}

	if info.Username == "" || info.Password == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrSetLoginAndPassword)
	}

	err := c.Validate(info)
	if err != nil {
		err := utils.GetValidationError(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	exists, err := h.Database.CheckUserExists(info.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	} else if exists {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrLoginAlreadyUse)
	}

	hashedPassword, err := auth.HashPassword(info.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	}

	err = h.Database.AddUser(info.Username, hashedPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	}

	token, err := auth.GenerateJWT(info.Username, h.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) Login(c echo.Context) error {
	info := new(models.User)

	if err := c.Bind(info); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrSetLoginAndPassword)
	}

	err := c.Validate(info)
	if err != nil {
		err := utils.GetValidationError(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	hashedPassword, err := h.Database.GetUser(info.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrUserDoesntExists)
	}

	if ok := auth.CheckPasswordHash(hashedPassword, info.Password); !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrLoginOrPassword)
	}

	token, err := auth.GenerateJWT(info.Username, h.JWTSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrTryLater)
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   0,
	})

	return c.Redirect(http.StatusSeeOther, "/")
}
