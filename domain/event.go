package domain

import "fmt"

type EventKind int

const (
	UnknownEvent    EventKind = 0
	InitializeEvent EventKind = 1
	DepositEvent    EventKind = 2
	WithdrawEvent   EventKind = 3
)

type Event struct {
	Kind      EventKind
	Money     Money
	EventId   int
	AccountId AccountId
}

func (e Event) String() string {
	return fmt.Sprintf("{kind=%v eventId=%d}", e.Kind, e.EventId)
}

func (k EventKind) String() string {
	switch k {
	case UnknownEvent:
		return "unknown"
	case InitializeEvent:
		return "initialize"
	case DepositEvent:
		return "deposit"
	case WithdrawEvent:
		return "withdraw"
	default:
		return "unknown"
	}
}
