package domain

import (
	"fmt"
)

type AccountId int64

type Account struct {
	state  accountState
	Events []Event
}

type accountState struct {
	accountId AccountId
	version   int
	money     Money
}

func CreateAccount(accountId AccountId) *Account {
	account := new(Account)
	account.initialize(accountId)
	return account
}

func RestoreAccount(events []Event) (*Account, error) {
	account := new(Account)
	for i := range events {
		err := account.appendEvent(events[i])
		if err != nil {
			return nil, fmt.Errorf("appendEvent: %w", err)
		}
	}
	return account, nil
}

func (s *accountState) applyEvent(event Event) error {
	switch event.Kind {
	case InitializeEvent:
		s.initialize(event)
	case DepositEvent:
		s.depositMoney(event)
	case WithdrawEvent:
		s.withdrawMoney(event)
	default:
		return fmt.Errorf("unknown event: %d", event.Kind)
	}
	return nil
}

func (s *Account) appendEvent(event Event) error {
	s.Events = append(s.Events, event)
	return s.state.applyEvent(event)
}

func (s *accountState) initialize(event Event) {
	s.accountId = event.AccountId
	s.version = 0
}

func (s *accountState) depositMoney(event Event) {
	s.money.Amount += event.Money.Amount
	s.version++
}

func (s *accountState) withdrawMoney(event Event) {
	s.money.Amount -= event.Money.Amount
	s.version++
}

func (s *Account) initialize(accountId AccountId) {
	s.appendEvent(Event{
		Kind:      InitializeEvent,
		Money:     Money{},
		AccountId: accountId,
	})
}

func (s *Account) DepositMoney(money Money) {
	s.appendEvent(Event{
		AccountId: s.state.accountId,
		Kind:      DepositEvent,
		Money:     money,
	})
}

func (s *Account) WithdrawMoney(money Money) error {
	if s.state.money.Amount < money.Amount {
		return fmt.Errorf("not enough money: %d", s.state.money.Amount)
	}
	s.appendEvent(Event{
		AccountId: s.state.accountId,
		Kind:      WithdrawEvent,
		Money:     money,
	})
	return nil
}

func (s *Account) Id() AccountId {
	return s.state.accountId
}

func (s *Account) Amount() Amount {
	return s.state.money.Amount
}

func (s *Account) Version() int {
	return s.state.version
}
