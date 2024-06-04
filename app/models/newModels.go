package models

import (
	"time"
)

type NewBank struct {
	GlobalIdentifier string
	Name             string
	Address          string
	CountryId        int
	BankTypeId       int
}

type NewBankClient struct {
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

type NewPolicyType struct {
	Id   int
	Code string
	Name string
}

type NewCountry struct {
	Id          int
	Name        string
	CountryCode string
}

type NewTransaction struct {
	Id                string
	OriginatorBankId  string
	BeneficiaryBankId string
	SenderId          string
	ReceiverId        string
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

type NewTransactionStatusHistory struct {
	TransactionId string
	StatusId      int
	Date          time.Time
}
type NewPolicy struct {
	Id                       int64
	PolicyTypeId             int
	TransactionTypeId        int
	PolicyEnforcingCountryId int
	OriginatingCountryId     int
	Parameters               string
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
