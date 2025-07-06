package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"urlshort/internal/database"
	"urlshort/internal/mocks"
	"urlshort/internal/models"
	"urlshort/internal/utils"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestDeleteLink_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/delete/testalias", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/delete/:alias")
	c.SetParamNames("alias")
	c.SetParamValues("testalias")

	user := models.UserToken{Username: "test-user", Ok: true}
	c.Set("username", user)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mocks.NewMockDatabase(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	mockDB.EXPECT().DeleteLink("testalias", "test-user").Return(true, nil)
	mockCache.EXPECT().DeleteCache("testalias").Return(nil)

	handler := NewHandlers(mockDB, mockCache, true, nil)

	assert.NoError(t, handler.DeleteLink(c))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteLink_Failed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/delete/testalias", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/delete/:alias")
	c.SetParamNames("alias")
	c.SetParamValues("testalias")

	user := models.UserToken{Username: "test-user", Ok: true}
	c.Set("username", user)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mocks.NewMockDatabase(ctrl)
	mockCache := mocks.NewMockCache(ctrl)

	mockDB.EXPECT().DeleteLink("testalias", "test-user").Return(false, nil)

	handler := NewHandlers(mockDB, mockCache, false, nil)

	err := handler.DeleteLink(c)
	assert.Error(t, err)
	assert.Equal(t, err, &echo.HTTPError{Message: ErrLinkNotFound, Code: http.StatusNotFound})
}

func TestShortUrl_Failed(t *testing.T) {

	e := echo.New()
	e.Validator = &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	f := make(url.Values)
	f.Set("link", "https://google.com")
	f.Set("alias", "123456789123456789")
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandlers(nil, nil, false, nil)

	err := handler.ShortUrl(c)
	assert.Error(t, err)
	assert.Condition(t, func() bool {
		return reflect.DeepEqual(err, &echo.HTTPError{Message: utils.ErrBadAlias, Code: http.StatusBadRequest}) ||
			reflect.DeepEqual(err, &echo.HTTPError{Message: utils.ErrBadLink, Code: http.StatusNotFound})
	}, "ожидаемая ошибка %v или %v, получено %v",
		&echo.HTTPError{Message: utils.ErrBadAlias, Code: http.StatusBadRequest},
		&echo.HTTPError{Message: utils.ErrBadLink, Code: http.StatusNotFound},
		err)
}

func TestShortUrl_Failed2(t *testing.T) {
	link := "https://google.com"
	alias := "123456"
	username := "test-user"

	e := echo.New()
	e.Validator = &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	f := make(url.Values)
	f.Set("link", link)
	f.Set("alias", alias)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := models.UserToken{Username: username, Ok: true}
	c.Set("username", user)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mocks.NewMockDatabase(ctrl)

	mockDB.EXPECT().InsertNew(models.ShortLink{Link: link, Alias: alias}, username).Return("", database.ErrAliasExists)

	handler := NewHandlers(mockDB, nil, false, nil)

	handler.ShortUrl(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, `{"error":"Ошибка алиаса! Такой алиас уже существует"}`+"\n", rec.Body.String())
}
