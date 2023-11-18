package usecase

import (
	"context"

	"github.com/waffleboot/ddd2/domain"
)

type DepositMoneyCommand struct {
	AccountId int
	Money     domain.Money
}

type WithdrawMoneyCommand struct {
	AccountId int
	Money     domain.Money
}

type CreateAccountUseCase interface {
	CreateAccount(ctx context.Context) (*domain.Account, error)
}

type DepositMoneyUseCase interface {
	DepositMoney(ctx context.Context, cmd DepositMoneyCommand) error
}

type WithdrawMoneyUseCase interface {
	WithdrawMoney(ctx context.Context, cmd WithdrawMoneyCommand) error
}
