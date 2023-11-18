package app

import (
	"context"
	"fmt"
	"sync/atomic"

	port "github.com/waffleboot/ddd2/app/port/database/account"
	"github.com/waffleboot/ddd2/app/port/usecase"
	"github.com/waffleboot/ddd2/domain"
)

var (
	_ usecase.GetAccountUseCase    = (*Service)(nil)
	_ usecase.CreateAccountUseCase = (*Service)(nil)
	_ usecase.DepositMoneyUseCase  = (*Service)(nil)
	_ usecase.WithdrawMoneyUseCase = (*Service)(nil)
)

type Service struct {
	repo   port.AccountRepo
	nextId int64
}

func New(repo port.AccountRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAccount(ctx context.Context) (*domain.Account, error) {
	accountId := atomic.AddInt64(&s.nextId, 1)

	account := domain.CreateAccount(domain.AccountId(accountId))

	err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("repo create: %w", err)
	}

	return account, nil
}

func (s *Service) GetAccount(ctx context.Context, accountId domain.AccountId) (*domain.Account, error) {
	account, err := s.getAccount(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) DepositMoney(ctx context.Context, cmd usecase.DepositMoneyCommand) error {
	account, err := s.getAccount(ctx, cmd.AccountId)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	prevVersion := account.Version()

	account.DepositMoney(cmd.Money)

	err = s.repo.Save(ctx, account, prevVersion)
	if err != nil {
		return fmt.Errorf("repo save: %w", err)
	}

	return nil
}

func (s *Service) WithdrawMoney(ctx context.Context, cmd usecase.WithdrawMoneyCommand) error {
	account, err := s.getAccount(ctx, cmd.AccountId)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	prevVersion := account.Version()

	err = account.WithdrawMoney(cmd.Money)
	if err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}

	err = s.repo.Save(ctx, account, prevVersion)
	if err != nil {
		return fmt.Errorf("repo save: %w", err)
	}

	return nil
}

func (s *Service) getAccount(ctx context.Context, accountId domain.AccountId) (*domain.Account, error) {
	events, err := s.repo.GetEvents(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}

	account, err := domain.RestoreAccount(events)
	if err != nil {
		return nil, fmt.Errorf("restore account: %w", err)
	}

	return account, nil
}
