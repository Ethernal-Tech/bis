package DB

import "time"

type BankEmployeeModel struct {
	Id       uint64
	Name     string
	Username string
	Password string
	BankId   uint64
	BankName string
}

type TransactionModel struct {
	Id              uint64
	OriginatorBank  string
	BeneficiaryBank string
	Sender          uint64
	Receiver        uint64
	Curency         string
	Amount          int
	Type            string
	Status          string
	StatusHistory   []StatusHistoryModel
	Policies        []string
}

type StatusHistoryModel struct {
	Date time.Time
	Name string
}
