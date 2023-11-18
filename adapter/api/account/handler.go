package account

import (
	"context"
	"fmt"

	"github.com/waffleboot/ddd2/app/port/usecase"
	"github.com/waffleboot/ddd2/domain"
)

type Handler struct {
	GetAccountUseCase    usecase.GetAccountUseCase
	CreateAccountUseCase usecase.CreateAccountUseCase
	DepositMoneyUseCase  usecase.DepositMoneyUseCase
	WithdrawMoneyUseCase usecase.WithdrawMoneyUseCase
}

func (s *Handler) GetAccount(ctx context.Context, accountId domain.AccountId) (*domain.Account, error) {
	return s.GetAccountUseCase.GetAccount(ctx, accountId)
}

func (s *Handler) CreateAccount(ctx context.Context) (domain.AccountId, error) {
	account, err := s.CreateAccountUseCase.CreateAccount(ctx)
	if err != nil {
		return 0, fmt.Errorf("create account: %w", err)
	}

	return account.Id(), nil
}

func (s *Handler) DepositAccount(ctx context.Context, accountId domain.AccountId, amount domain.Amount) error {
	cmd, err := usecase.NewDepositMoneyCommand(accountId, amount)
	if err != nil {
		return fmt.Errorf("new deposit money command: %w", err)
	}
	return s.DepositMoneyUseCase.DepositMoney(ctx, cmd)
}

func (s *Handler) WithdrawAccount(ctx context.Context, accountId domain.AccountId, amount domain.Amount) error {
	cmd, err := usecase.NewWithdrawMoneyCommand(accountId, amount)
	if err != nil {
		return fmt.Errorf("new deposit money command: %w", err)
	}
	return s.WithdrawMoneyUseCase.WithdrawMoney(ctx, cmd)
}
