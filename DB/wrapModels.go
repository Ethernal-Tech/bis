package DB

import "time"

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
}

type StatusHistoryModel struct {
	Date time.Time
	Name string
}
