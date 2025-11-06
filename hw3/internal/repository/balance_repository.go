package repository

import "homework3/internal/domain"

type BalanceRepository interface {
	GetBalance(userID string) (*domain.Balance, error)
	CreateBalance(userID string) (*domain.Balance, error)
	UpdateBalance(balance *domain.Balance) error
	GetOrCreateBalance(userID string) (*domain.Balance, error)
}
