package main

import (
	"errors"
	"sync"
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

var errNotFound = errors.New("not found")

type store struct {
	mu     sync.Mutex
	nextID int
	todos  map[int]*Todo
	order  []int
}

func newStore() *store {
	s := &store{todos: map[int]*Todo{}}
	s.create("Turbo Drive を試す")
	s.create("Turbo Frames に着手する")
	return s
}

func (s *store) create(title string) *Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	t := &Todo{ID: s.nextID, Title: title}
	s.todos[t.ID] = t
	s.order = append(s.order, t.ID)
	return t
}

func (s *store) list() []*Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Todo, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.todos[id])
	}
	return out
}

func (s *store) get(id int) (*Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	return t, ok
}

func (s *store) update(id int, title string, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.todos[id]
	if !ok {
		return errNotFound
	}
	t.Title = title
	t.Done = done
	return nil
}

func (s *store) delete(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return
	}
	delete(s.todos, id)
	for i, v := range s.order {
		if v == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}
