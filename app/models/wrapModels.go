package models

import "time"

type BankEmployeeModel struct {
	Name     string
	Username string
	Password string
	BankId   string
	BankName string
	Country  string
}

type TransactionModel struct {
	Id                       string
	OriginatorBank           string
	BeneficiaryBank          string
	SenderGlobalIdentifier   string
	ReceiverGlobalIdentifier string
	SenderName               string
	ReceiverName             string
	Currency                 string
	Amount                   int
	LoanId                   int
	Type                     string
	TypeCode                 string
	TypeId                   int
	Status                   string
	StatusHistory            []StatusHistoryModel
	Policies                 []string
	OriginatorBankCountryId  string
}

type StatusHistoryModel struct {
	Date       time.Time
	DateString string
	Name       string
}

type PolicyModel struct {
	Id              uint64
	Country         string
	CountryId       int
	Code            string
	Name            string
	Amount          uint64
	Checklist       string
	Parameter       string
	TransactionType string
}

type InteractiveComplianceCheckProofRequest struct {
	Value             string `json:"value"`
	ComplianceCheckID string `json:"compliance_check_id"`
	PolicyID          string `json:"policy_id"`
}

type BankModel struct {
	Id      string
	Name    string
	Country string
}

type SearchModel struct {
	Value   string `json:"value"`
	From    string `json:"from"`
	To      string `json:"to"`
	StateId string `json:StateId`
}
