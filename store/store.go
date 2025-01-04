package store

type Store interface {
	CreateItem(key string, value string) error
	GetItem(key string) (string, error)
	UpdateItem(key string, value string) error
	DeleteItem(key string) error
}
