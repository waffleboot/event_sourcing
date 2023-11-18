package database

import (
	"context"

	"github.com/waffleboot/ddd2/domain"
)

type AccountRepo interface {
	GetEvents(ctx context.Context, accountId int) ([]domain.Event, error)
	Save(ctx context.Context, account *domain.Account) error
}
