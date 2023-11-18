package app

import (
	"context"
	"fmt"

	"github.com/waffleboot/ddd2/app/port/database"
	"github.com/waffleboot/ddd2/app/port/usecase"
	"github.com/waffleboot/ddd2/domain"
)

var (
	_ usecase.CreateAccountUseCase = (*Service)(nil)
	_ usecase.DepositMoneyUseCase  = (*Service)(nil)
	_ usecase.WithdrawMoneyUseCase = (*Service)(nil)
)

type Service struct {
	repo database.AccountRepo
}

func (*Service) CreateAccount(ctx context.Context) (*domain.Account, error) {
	return domain.CreateAccount(), nil
}

func (s *Service) DepositMoney(ctx context.Context, cmd usecase.DepositMoneyCommand) error {
	events, err := s.repo.GetEvents(ctx, cmd.AccountId)
	if err != nil {
		return fmt.Errorf("get events: %w", err)
	}

	account, err := domain.RestoreAccount(events)
	if err != nil {
		return fmt.Errorf("restore account: %w", err)
	}

	account.DepositMoney(cmd.Money)

	return nil
}

func (s *Service) WithdrawMoney(ctx context.Context, cmd usecase.WithdrawMoneyCommand) error {
	events, err := s.repo.GetEvents(ctx, cmd.AccountId)
	if err != nil {
		return fmt.Errorf("get events: %w", err)
	}

	account, err := domain.RestoreAccount(events)
	if err != nil {
		return fmt.Errorf("restore account: %w", err)
	}

	err = account.WithdrawMoney(cmd.Money)
	if err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}

	return nil
}
