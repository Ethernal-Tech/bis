package models

import (
	"time"
)

type NewBank struct {
	GlobalIdentifier string
	Name             string
	Address          string
	CountryId        uint
	BankTypeId       uint
}

type NewBankClient struct {
	Id               uint
	GlobalIdentifier string
	Name             string
	Address          string
	BankId           string
}

type NewBankEmployee struct {
	Id       string
	Name     string
	Username string
	Password string
	BankId   string
}

type NewBankType struct {
	Id   int
	Name string
}

type NewCountry struct {
	Id   int
	Name string
	Code string
}

type NewTransaction struct {
	Id                string
	OriginatorBankId  string
	BeneficiaryBankId string
	SenderId          uint
	ReceiverId        uint
	Currency          string
	Amount            int
	TransactionTypeId int
	LoanId            int
}

type NewTransactionType struct {
	Id   int
	Code string
	Name string
}

type NewStatus struct {
	Id   int
	Name string
}

type NewTransactionHistory struct {
	TransactionId string
	StatusId      int
	Date          time.Time
}

type NewPolicyType struct {
	Id   int
	Code string
	Name string
}

type NewPolicy struct {
	Id                       int
	PolicyTypeId             int
	TransactionTypeId        int
	PolicyEnforcingCountryId int
	OriginatingCountryId     int
	Parameters               string
}

type NewPolicyModel struct {
	Policy     NewPolicy
	PolicyType NewPolicyType
}

type NewTransactionPolicy struct {
	TransactionId string
	PolicyId      int64
	Status        int
}

type NewTransactionProof struct {
	Id            uint64
	TransactionId uint64
	Proof         string
}

// TODO: Remove ASSET statuses
// const (
// 	Created         int = 0
// 	PoliciesApplied int = 1
// 	ProofRequested  int = 2
// 	ProofReceived   int = 3
// 	ProofInvalid    int = 4
// 	AssetSent       int = 5
// 	AssetReceived   int = 6
// 	Canceled        int = 7
// )

// const (
// 	CommercialBank = iota + 1
// 	CentralBank
// )
