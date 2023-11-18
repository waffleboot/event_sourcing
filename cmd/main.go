package main

import (
	"context"
	"fmt"

	"github.com/waffleboot/ddd2/adapter/api"
	accountApi "github.com/waffleboot/ddd2/adapter/api/account"
	"github.com/waffleboot/ddd2/adapter/database/account"
	"github.com/waffleboot/ddd2/app"
	"github.com/waffleboot/ddd2/app/port/usecase"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func run(ctx context.Context) error {
	accountRepo := account.New()
	service := app.New(accountRepo)

	account, err := service.CreateAccount(ctx)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	cmd, err := usecase.NewDepositMoneyCommand(account.Id(), 500)
	if err != nil {
		return fmt.Errorf("new deposit money command: %w", err)
	}

	err = service.DepositMoney(ctx, cmd)
	if err != nil {
		return fmt.Errorf("deposit money: %w", err)
	}

	cmd2, err := usecase.NewWithdrawMoneyCommand(account.Id(), 50)
	if err != nil {
		return fmt.Errorf("new deposit money command: %w", err)
	}

	err = service.WithdrawMoney(ctx, cmd2)
	if err != nil {
		return fmt.Errorf("deposit money: %w", err)
	}

	err = service.WithdrawMoney(ctx, cmd2)
	if err != nil {
		return fmt.Errorf("deposit money: %w", err)
	}

	handler := &accountApi.Handler{
		GetAccountUseCase:    service,
		CreateAccountUseCase: service,
		DepositMoneyUseCase:  service,
		WithdrawMoneyUseCase: service,
	}

	err = api.Start(handler)
	if err != nil {
		return fmt.Errorf("start api: %w", err)
	}

	return nil
}
