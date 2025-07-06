package utils

import (
	"github.com/go-playground/validator/v10"
)

var (
	ErrBadLink  = "Плохая ссылка! Максимальная длина 150 символов"
	ErrBadAlias = "Плохой алиас! Максимальная длина 12 символов"
	ErrUnknown  = "Неизвестная ошибка"
	ErrUsername = "Логин должен содержать только латиницу и цифры. Максимальная длина 12 символов."
	ErrPassword = "Пожалуйста, используйте пароль короче 128 символов :)"
)

func GetValidationError(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			switch e.Field() {
			case "Link":
				return ErrBadLink
			case "Alias":
				return ErrBadAlias
			case "Username":
				return ErrUsername
			case "Password":
				return ErrPassword
			}
		}
	}
	return ErrUnknown
}
