package store

import (
	"errors"
	"sync"

	"github.com/asdlc-repos/testingnewtaskpage3979/user-service/internal/models"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidDays        = errors.New("days must be greater than zero")
)

type Store struct {
	mu       sync.RWMutex
	users    map[string]*models.User
	balances map[string]*models.LeaveBalance
}

func New() *Store {
	s := &Store{
		users:    make(map[string]*models.User),
		balances: make(map[string]*models.LeaveBalance),
	}
	s.seed()
	return s
}

func (s *Store) seed() {
	users := []*models.User{
		{Id: "u1", Name: "Alice Johnson", Email: "alice@example.com", Role: "manager"},
		{Id: "u2", Name: "Bob Smith", Email: "bob@example.com", Role: "manager"},
		{Id: "u3", Name: "Carol White", Email: "carol@example.com", Role: "employee", ManagerId: "u1"},
		{Id: "u4", Name: "David Brown", Email: "david@example.com", Role: "employee", ManagerId: "u1"},
		{Id: "u5", Name: "Eva Martinez", Email: "eva@example.com", Role: "employee", ManagerId: "u1"},
		{Id: "u6", Name: "Frank Lee", Email: "frank@example.com", Role: "employee", ManagerId: "u2"},
		{Id: "u7", Name: "Grace Kim", Email: "grace@example.com", Role: "employee", ManagerId: "u2"},
	}
	for _, u := range users {
		s.users[u.Id] = u
		s.balances[u.Id] = &models.LeaveBalance{
			UserId:      u.Id,
			Entitlement: 20.0,
			Used:        0.0,
			Remaining:   20.0,
		}
	}
}

func (s *Store) GetUser(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *u
	return &cp, nil
}

func (s *Store) GetBalance(userId string) (*models.LeaveBalance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.users[userId]; !ok {
		return nil, ErrNotFound
	}
	b, ok := s.balances[userId]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *b
	return &cp, nil
}

func (s *Store) DeductBalance(userId string, days float64) (*models.LeaveBalance, error) {
	if days <= 0 {
		return nil, ErrInvalidDays
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userId]; !ok {
		return nil, ErrNotFound
	}
	b, ok := s.balances[userId]
	if !ok {
		return nil, ErrNotFound
	}
	if days > b.Remaining {
		return nil, ErrInsufficientBalance
	}
	b.Used += days
	b.Remaining = b.Entitlement - b.Used
	cp := *b
	return &cp, nil
}

func (s *Store) GetDirectReports(managerId string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	mgr, ok := s.users[managerId]
	if !ok {
		return nil, ErrNotFound
	}
	_ = mgr
	var reports []string
	for _, u := range s.users {
		if u.ManagerId == managerId {
			reports = append(reports, u.Id)
		}
	}
	return reports, nil
}
