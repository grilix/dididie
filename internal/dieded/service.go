package dieded

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Service interface {
	CreateProfile(ctx context.Context, p ProfileForm) (*Profile, error)
	GetProfile(ctx context.Context, id int) (*Profile, error)
	QueryProfile(ctx context.Context, q Query) (*Profile, error)
	DieProfile(ctx context.Context, id int) (*Profile, error)
	DeleteProfile(ctx context.Context, id int) error
}

type Query struct {
	Name string `json:"name,omitempty"`
}

type DieRequest struct {
	ID int `json:"id"`
}

type Profile struct {
	ID     int        `json:"id"`
	Name   string     `json:"name,omitempty"`
	DiedAt *time.Time `json:"died_at,omitempty"`
}

type ProfileForm struct {
	Name string `json:"name"`
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyDied   = errors.New("already died")
	ErrDuplicateName = errors.New("duplicate name")
)

type inmemService struct {
	mtx    sync.RWMutex
	m      map[int]Profile
	nextID int
}

func NewInmemService() Service {
	return &inmemService{
		m: map[int]Profile{},
	}
}

func (s *inmemService) CreateProfile(ctx context.Context, f ProfileForm) (*Profile, error) {
	_, err := s.QueryProfile(ctx, Query{Name: f.Name})
	if err == nil {
		return nil, ErrDuplicateName
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	now := time.Now()
	p := Profile{
		ID:     s.nextID,
		Name:   f.Name,
		DiedAt: &now,
	}
	s.nextID++
	s.m[p.ID] = p

	return &p, nil
}

func (s *inmemService) GetProfile(ctx context.Context, id int) (*Profile, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &p, nil
}

func (s *inmemService) QueryProfile(ctx context.Context, q Query) (*Profile, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, p := range s.m {
		if p.Name == q.Name {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *inmemService) DieProfile(ctx context.Context, id int) (*Profile, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	p, ok := s.m[id]
	if !ok {
		return nil, ErrNotFound
	}
	if p.DiedAt != nil {
		diedAt := p.DiedAt.Format("1999-12-31")
		if diedAt == time.Now().Format("1999-12-31") {
			return &p, ErrAlreadyDied
		}
	}
	now := time.Now()
	p.DiedAt = &now
	s.m[id] = p
	return &p, nil
}

func (s *inmemService) DeleteProfile(ctx context.Context, id int) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}
