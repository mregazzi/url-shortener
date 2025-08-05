package storage

type Store interface {
	Save(code, url string) error
	Get(code string) (string, bool, error)
}
