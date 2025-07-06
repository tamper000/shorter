//go:generate mockgen -source=database.go -destination=../mocks/database_mock.go --package mocks

package interfaces

import "urlshort/internal/models"

type Database interface {
	InsertNew(m models.ShortLink, username string) (string, error)
	GetByAlias(alias string) (string, error)
	AddClick(alias string) error
	CheckUserExists(username string) (bool, error)
	GetUser(username string) (string, error)
	AddUser(username, password string) error
	GetLinksByUser(username string) ([]models.Link, error)
	DeleteLink(alias, username string) (bool, error)
	CloseDatabase()
}
