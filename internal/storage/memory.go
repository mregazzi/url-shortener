package storage

import "sync"

var _ Store = (*InMemoryStore)(nil)

type InMemoryStore struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (s *InMemoryStore) Save(code, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[code] = url
	return nil
}

func (s *InMemoryStore) Get(code string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[code]
	return url, ok, nil
}
