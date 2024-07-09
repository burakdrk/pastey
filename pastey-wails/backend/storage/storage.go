package storage

type Storage interface {
	Save(key string, value string)
	Get(key string) (string, error)
	Delete(key string)
	Exists(key string) bool
	Close()
}
