package account

import (
	"context"
	"errors"

	"github.com/waffleboot/ddd2/domain"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type AccountRepo interface {
	// throw ErrConflict
	Create(ctx context.Context, account *domain.Account) error
	// throw ErrorNotFound, ErrConflict
	Save(ctx context.Context, account *domain.Account, prevVersion int) error
	// throw ErrorNotFound
	GetEvents(ctx context.Context, accountId domain.AccountId) ([]domain.Event, error)
}
