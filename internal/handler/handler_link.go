package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"urlshort/internal/models"
	"urlshort/internal/utils"
)

func (h *Handler) ShortUrl(c echo.Context) error {
	info := new(models.ShortLink)

	if err := c.Bind(info); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrBadRequest)
	}

	err := c.Validate(info)
	if err != nil {
		err := utils.GetValidationError(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var username string
	user := c.Get("username").(models.UserToken)
	if user.Ok {
		username = user.Username
	}

	alias, err := h.Database.InsertNew(*info, username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseShortLink{Error: "Ошибка алиаса! " + utils.FirstUpper(err.Error())})
	}

	if h.IsCacheEnabled {
		h.Cache.AddCache(alias, info.Link)
	}

	return c.JSON(http.StatusOK, models.ResponseShortLink{
		Alias: alias,
	})
}

func (h *Handler) RedirectHandler(c echo.Context) error {
	alias := c.Param("alias")

	if h.IsCacheEnabled {
		link, err := h.Cache.GetCache(alias)
		if err == nil {
			h.Database.AddClick(alias)
			return c.Redirect(http.StatusPermanentRedirect, link)
		}
	}

	link, err := h.Database.GetByAlias(alias)
	if err != nil {
		return c.Redirect(http.StatusPermanentRedirect, "/")
	}
	h.Database.AddClick(alias)
	if h.IsCacheEnabled {
		h.Cache.AddCache(alias, link)
	}

	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

	return c.Redirect(http.StatusPermanentRedirect, link)
}

func (h *Handler) DeleteLink(c echo.Context) error {
	user := c.Get("username").(models.UserToken)
	if !user.Ok {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrUnauthorized)
	}

	alias := c.Param("alias")
	isAffected, err := h.Database.DeleteLink(alias, user.Username)
	if err != nil {
		return err
	} else if !isAffected {
		return echo.NewHTTPError(http.StatusNotFound, ErrLinkNotFound)
	}

	if h.IsCacheEnabled {
		h.Cache.DeleteCache(alias)
	}
	return c.NoContent(http.StatusOK)
}
