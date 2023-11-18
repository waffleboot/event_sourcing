package account

import (
	"context"
	"fmt"
	"sync"

	port "github.com/waffleboot/ddd2/app/port/database/account"
	"github.com/waffleboot/ddd2/domain"
)

var _ port.AccountRepo = (*Repo)(nil)

type event struct {
	kind      int
	eventId   int
	accountId domain.AccountId
	amount    domain.Amount
}

type Repo struct {
	events   map[domain.AccountId][]event
	versions map[domain.AccountId]int
	mu       sync.RWMutex
}

func New() *Repo {
	return &Repo{
		events:   make(map[domain.AccountId][]event),
		versions: make(map[domain.AccountId]int)}
}

func (s *Repo) Create(ctx context.Context, account *domain.Account) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.versions[account.Id()]
	if ok {
		return fmt.Errorf("account already exists: %w", port.ErrConflict)
	}

	s.versions[account.Id()] = account.Version()

	err := s.saveEvents(ctx, account, false)
	if err != nil {
		return fmt.Errorf("save events: %w", err)
	}

	return nil
}

func (s *Repo) GetEvents(ctx context.Context, accountId domain.AccountId) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[accountId]
	if !ok {
		return nil, fmt.Errorf("account not found: %w", port.ErrNotFound)
	}

	domainEvents := make([]domain.Event, 0, len(events))
	for i := range events {
		domainEvents = append(domainEvents, events[i].domain())
	}

	return domainEvents, nil
}

func (s *Repo) Save(ctx context.Context, account *domain.Account, prevVersion int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	version, ok := s.versions[account.Id()]
	if !ok {
		return fmt.Errorf("account not found: %w", port.ErrNotFound)
	}

	if version != prevVersion {
		return fmt.Errorf("inconsistent version: current=%d prev=%d %w", version, prevVersion, port.ErrConflict)
	}

	s.versions[account.Id()] = account.Version()

	err := s.saveEvents(ctx, account, true)
	if err != nil {
		return fmt.Errorf("save events: %w", err)
	}

	return nil
}

func (s *Repo) saveEvents(ctx context.Context, account *domain.Account, check bool) error {
	events, ok := s.events[account.Id()]
	if check && !ok {
		return fmt.Errorf("events not found: %w", port.ErrNotFound)
	}

	for i := range account.Events {
		if account.Events[i].EventId == 0 {
			var event event
			event.bind(account.Events[i])
			event.eventId = len(events) + 1
			events = append(events, event)
		}
	}

	s.events[account.Id()] = events

	return nil
}

func (s *event) domain() domain.Event {
	return domain.Event{
		Kind: domain.EventKind(s.kind),
		Money: domain.Money{
			Amount: s.amount,
		},
		EventId:   s.eventId,
		AccountId: s.accountId,
	}
}

func (s *event) bind(e domain.Event) {
	s.kind = int(e.Kind)
	s.amount = e.Money.Amount
	s.eventId = e.EventId
	s.accountId = e.AccountId
}
