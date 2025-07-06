//go:generate mockgen -source=cache.go -destination=../mocks/cache_mock.go --package mocks

package interfaces

type Cache interface {
	AddCache(alias, link string) error
	GetCache(alias string) (string, error)
	DeleteCache(key string) error
	CloseCache()
}
