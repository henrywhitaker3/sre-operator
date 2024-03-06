package store

import (
	"context"
	"errors"
	"sync"
)

const (
	WEBHOOK = "webhook"
)

var (
	ErrInvalidType    = errors.New("invalid type")
	ErrInvalidSubFunc = errors.New("nil subscription function provided")
	ErrUnknownTrigger = errors.New("tried to subscribe to an unknown trigger")
	ErrNotSubscribed  = errors.New("not subscribe to that tirgger")
)

type StoreSubscriber func(ctx context.Context) error

type Subscription struct {
	Name string
	Do   StoreSubscriber
}

type Store struct {
	// A map of items keyed by the kind of trigger e.g.
	// [
	//     WEBHOOK: [
	//         "demo-webhook": [
	//             {Name: "demo-rollout", Do: func() {}}
	//         ]
	//     ]
	// ]
	items map[string]map[string][]Subscription
	mu    *sync.Mutex
}

func NewStore() *Store {
	return &Store{
		items: map[string]map[string][]Subscription{
			WEBHOOK: make(map[string][]Subscription),
		},
		mu: &sync.Mutex{},
	}
}

// Register a new trigger into the store
func (s *Store) Register(kind string, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[kind]; !ok {
		return ErrInvalidType
	}

	if _, ok := s.items[kind][id]; !ok {
		s.items[kind][id] = []Subscription{}
	}
	return nil
}

func (s *Store) Subscibe(kind string, id string, sub Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[kind]; !ok {
		return ErrInvalidType
	}
	if _, ok := s.items[kind][id]; !ok {
		return ErrUnknownTrigger
	}
	if sub.Do == nil {
		return ErrInvalidSubFunc
	}
	s.items[kind][id] = append(s.items[kind][id], sub)
	return nil
}

func (s Store) Unsubscribe(kind string, id string, sub Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[kind]; !ok {
		return ErrInvalidType
	}
	if _, ok := s.items[kind][id]; !ok {
		return ErrUnknownTrigger
	}

	for i, ss := range s.items[kind][id] {
		if ss.Name == sub.Name {
			s.items[kind][id] = append(s.items[kind][id][:i], s.items[kind][id][i+1:]...)
			return nil
		}
	}
	return ErrNotSubscribed
}

func (s *Store) Get(kind string, id string) ([]Subscription, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[kind]; !ok {
		return nil, ErrInvalidType
	}
	if _, ok := s.items[kind][id]; !ok {
		return nil, ErrUnknownTrigger
	}
	return s.items[kind][id], nil
}

func (s *Store) Drop(kind string, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[kind]; !ok {
		return ErrInvalidType
	}
	if _, ok := s.items[kind][id]; !ok {
		return nil
	}
	delete(s.items[kind], id)
	return nil
}
