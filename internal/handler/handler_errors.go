package handler

var (
	ErrBadRequest   = "Плохие данные. Перепроверьте все поля ввода."
	ErrUnauthorized = "Вы не авторизованы."
	ErrTryLater     = "Попробуйте позже."
	ErrLinkNotFound = "Такой ссылки не существует."
)

var (
	ErrSetLoginAndPassword = "Укажите логин или пароль."
	ErrLoginAlreadyUse     = "Такой логин уже используется."
	ErrUserDoesntExists    = "Такого пользователя не существует."
	ErrLoginOrPassword     = "Проверьте логин и пароль."
)
