package domain

import (
	"fmt"
)

type Account struct {
	State  AccountState
	Events []Event
}

type AccountState struct {
	AccountId int
	Version   int
	Money     Money
}

func CreateAccount() *Account {
	account := new(Account)
	account.initialize()
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

func (s *AccountState) applyEvent(event Event) error {
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
	return s.State.applyEvent(event)
}

func (s *AccountState) initialize(event Event) {
	s.AccountId = event.AccountId
	s.Version = 0
}

func (s *AccountState) depositMoney(event Event) {
	s.Money.Amount += event.Money.Amount
	s.Version++
}

func (s *AccountState) withdrawMoney(event Event) {
	s.Money.Amount -= event.Money.Amount
	s.Version++
}

func (s *Account) initialize() {
	s.appendEvent(Event{
		Kind:      InitializeEvent,
		Money:     Money{},
		AccountId: 0,
	})
}

func (s *Account) DepositMoney(money Money) {
	s.appendEvent(Event{
		AccountId: s.State.AccountId,
		Kind:      DepositEvent,
		Money:     money,
	})
}

func (s *Account) WithdrawMoney(money Money) error {
	if s.State.Money.Amount < money.Amount {
		return fmt.Errorf("not enough money: ", s.State.Money.Amount)
	}
	s.appendEvent(Event{
		AccountId: s.State.AccountId,
		Kind:      WithdrawEvent,
		Money:     money,
	})
	return nil
}
