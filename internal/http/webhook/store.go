package webhook

import (
	"context"
	"sync"
)

type StoreSubscriber func(ctx context.Context) error

type Store struct {
	hooks map[string]map[string]StoreSubscriber
	mu    *sync.Mutex
}

func NewStore() *Store {
	return &Store{
		hooks: make(map[string]map[string]StoreSubscriber),
		mu:    &sync.Mutex{},
	}
}

func (s *Store) Store(hook string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hooks[hook]; !ok {
		s.hooks[hook] = make(map[string]StoreSubscriber)
	}
}

func (s *Store) StoreFunc(hook string, name string, f StoreSubscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hooks[hook]; !ok {
		s.hooks[hook] = make(map[string]StoreSubscriber)
	}

	if _, ok := s.hooks[hook][name]; ok {
		return
	}

	if f == nil {
		return
	}

	s.hooks[hook][name] = f
}

func (s Store) Drop(hook string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hooks[hook]; !ok {
		return
	}

	delete(s.hooks, hook)
}

func (s *Store) DropFunc(hook string, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.hooks[hook]; !ok {
		return nil
	}

	if _, ok := s.hooks[hook][name]; !ok {
		return nil
	}

	delete(s.hooks[hook], name)
	return nil
}

func (s *Store) Get(hook string) (map[string]StoreSubscriber, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cb, ok := s.hooks[hook]; ok {
		return cb, true
	}
	return nil, false
}
