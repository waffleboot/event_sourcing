package usecase

import (
	"context"

	"github.com/waffleboot/ddd2/domain"
)

type DepositMoneyCommand struct {
	AccountId domain.AccountId
	Money     domain.Money
}

type WithdrawMoneyCommand struct {
	AccountId domain.AccountId
	Money     domain.Money
}

func NewDepositMoneyCommand(accountId domain.AccountId, amount domain.Amount) (DepositMoneyCommand, error) {
	return DepositMoneyCommand{
		AccountId: accountId,
		Money: domain.Money{
			Amount: amount,
		},
	}, nil
}

func NewWithdrawMoneyCommand(accountId domain.AccountId, amount domain.Amount) (WithdrawMoneyCommand, error) {
	return WithdrawMoneyCommand{
		AccountId: accountId,
		Money: domain.Money{
			Amount: amount,
		},
	}, nil
}

type CreateAccountUseCase interface {
	CreateAccount(ctx context.Context) (*domain.Account, error)
}

type GetAccountUseCase interface {
	GetAccount(ctx context.Context, accountId domain.AccountId) (*domain.Account, error)
}

type DepositMoneyUseCase interface {
	DepositMoney(ctx context.Context, cmd DepositMoneyCommand) error
}

type WithdrawMoneyUseCase interface {
	WithdrawMoney(ctx context.Context, cmd WithdrawMoneyCommand) error
}
