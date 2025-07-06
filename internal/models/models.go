package models

import (
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ShortLink struct {
	Link  string `param:"link" query:"link" form:"link" json:"link" xml:"link" validate:"required,url,max=150"`
	Alias string `param:"alias" query:"alias" form:"alias" json:"alias" xml:"alias" validate:"alphanum,max=12"`
}

type ResponseShortLink struct {
	Alias string `json:"alias,omitempty"`
	Error string `json:"error,omitempty"`
}

type User struct {
	Username string `json:"username" validate:"required,alphanum,max=12"`
	Password string `json:"password" validate:"required,max=128"`
}

func (u User) Validate() (bool, string) {
	if err := validate.Var(u.Username, "required,alphanum,max=12"); err != nil {
		return false, url.QueryEscape("Логин должен содержать только латиницу и цифры. Максимальная длина 12 символов.")
	}

	if err := validate.Var(u.Password, "required,max=128"); err != nil {
		return false, url.QueryEscape("Пожалуйста, используйте пароль короче 128 символов :)")
	}

	return true, ""
}

type UserToken struct {
	Username string
	Ok       bool
}

type Link struct {
	Alias    string
	Original string
	Clicks   int
}

type AccessToken struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
