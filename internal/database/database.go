package database

import (
	"errors"
)

var (
	ErrFailedGenerate = errors.New("не удалось сгенерировать алиас")
	ErrAliasExists    = errors.New("такой алиас уже существует")
	ErrFailedCreate   = errors.New("не удалось создать короткую ссылку")
)
