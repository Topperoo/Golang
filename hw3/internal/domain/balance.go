package domain

import (
	"errors"
	"sync"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrNegativeAmount = errors.New("amount cannot be negative")
	ErrSelfTransfer = errors.New("cannot transfer to yourself")
	ErrUserNotFound = errors.New("user not found")
)

type Balance struct {
	UserID string
	Amount float64
	mu sync.RWMutex
}

func NewBalance(userID string) *Balance {
	return &Balance{UserID: userID}
}

func (b *Balance) Credit(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.Amount += amount
	return nil
}

func (b *Balance) Debit(amount float64) error {
	if amount < 0 {
		return ErrNegativeAmount
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Amount < amount {
		return ErrInsufficientBalance
	}

	b.Amount -= amount
	return nil
}

func (b *Balance) GetAmount() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Amount
}
