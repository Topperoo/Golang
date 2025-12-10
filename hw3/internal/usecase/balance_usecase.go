package usecase

import (
	"homework3/internal/domain"
	"homework3/internal/repository"
	"sync"
)

type BalanceUseCase struct {
	repo repository.BalanceRepository
	mu sync.Mutex
}

func NewBalanceUseCase(repo repository.BalanceRepository) *BalanceUseCase {
	return &BalanceUseCase{repo: repo}
}

func (uc *BalanceUseCase) CreditBalance(userID string, amount float64) error {
	if amount <= 0 {
		return domain.ErrNegativeAmount
	}

	balance, err := uc.repo.GetOrCreateBalance(userID)
	if err != nil {
		return err
	}

	if err := balance.Credit(amount); err != nil {
		return err
	}

	return uc.repo.UpdateBalance(balance)
}

func (uc *BalanceUseCase) TransferBalance(fromUserID, toUserID string, amount float64) error {
	if amount <= 0 {
		return domain.ErrNegativeAmount
	}

	if fromUserID == toUserID {
		return domain.ErrSelfTransfer
	}

	uc.mu.Lock()
	defer uc.mu.Unlock()

	fromBalance, err := uc.repo.GetBalance(fromUserID)
	if err != nil {
		return err
	}

	toBalance, err := uc.repo.GetOrCreateBalance(toUserID)
	if err != nil {
		return err
	}

	if err := fromBalance.Debit(amount); err != nil {
		return err
	}

	if err := toBalance.Credit(amount); err != nil {
		_ = fromBalance.Credit(amount)
		return err
	}

	if err := uc.repo.UpdateBalance(fromBalance); err != nil {
		return err
	}

	return uc.repo.UpdateBalance(toBalance)
}

func (uc *BalanceUseCase) GetBalance(userID string) (float64, error) {
	balance, err := uc.repo.GetOrCreateBalance(userID)
	if err != nil {
		return 0, err
	}

	return balance.GetAmount(), nil
}
