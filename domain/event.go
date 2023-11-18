package domain

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
	AccountId int
}
