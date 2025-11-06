package memory

import (
	"homework3/internal/domain"
	"homework3/internal/repository"
	"sync"
)

type BalanceMemoryRepository struct {
	balances map[string]*domain.Balance
	mu       sync.RWMutex
}

func NewBalanceMemoryRepository() repository.BalanceRepository {
	return &BalanceMemoryRepository{
		balances: make(map[string]*domain.Balance),
	}
}

func (r *BalanceMemoryRepository) GetBalance(userID string) (*domain.Balance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	balance, exists := r.balances[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return balance, nil
}

func (r *BalanceMemoryRepository) CreateBalance(userID string) (*domain.Balance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.balances[userID]; exists {
		return r.balances[userID], nil
	}

	balance := domain.NewBalance(userID)
	r.balances[userID] = balance
	return balance, nil
}

func (r *BalanceMemoryRepository) UpdateBalance(balance *domain.Balance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.balances[balance.UserID] = balance
	return nil
}

func (r *BalanceMemoryRepository) GetOrCreateBalance(userID string) (*domain.Balance, error) {
	balance, err := r.GetBalance(userID)
	if err == domain.ErrUserNotFound {
		return r.CreateBalance(userID)
	}
	return balance, err
}
