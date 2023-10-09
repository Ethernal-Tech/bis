package DB

import (
	"time"
)

type Bank struct {
	Id               uint64
	GlobalIdentifier string
	Name             string
	Address          string
	Country          string
}

type BankClient struct {
	Id               uint64
	GlobalIdentifier string
	Name             string
	Address          string
	BankId           uint64
}

type BankEmployee struct {
	Id       uint64
	Name     string
	Username string
	Password string
	BankId   uint64
}

type Transaction struct {
	Id              uint64
	OriginatorBank  uint64
	BeneficiaryBank uint64
	Sender          uint64
	Receiver        uint64
	Date            time.Time
	Curency         string
	Amount          int
	TypeId          int
}

type Type struct {
	Id   int
	Name string
}

type Status struct {
	Id   int
	Name string
}

type TransactionStatusHistory struct {
	Id            int
	TransactionId uint64
	StatusId      int
	Date          time.Time
}

type Policy struct {
	Id   int
	Name string
}

type TransactionPolicy struct {
	TransactionId uint64
	PolicyId      int
}

type TransactionProof struct {
	Id            uint64
	TransactionId uint64
	Proof         string
}
