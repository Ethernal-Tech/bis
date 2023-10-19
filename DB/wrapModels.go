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
	Id                        uint64
	OriginatorBank            string
	BeneficiaryBank           string
	SenderGlobalIdentifier    string
	ReceiverGlobalIdedntifier string
	SenderName                string
	ReceiverName              string
	Currency                  string
	Amount                    int
	Type                      string
	Status                    string
	StatusHistory             []StatusHistoryModel
	Policies                  []string
}

type StatusHistoryModel struct {
	Date       time.Time
	DateString string
	Name       string
}

type PolicyModel struct {
	Country string
	Amount  uint64
	Name    string
}

type TransactionProofRequest struct {
	Value 		  string
	TransactionId uint64
}