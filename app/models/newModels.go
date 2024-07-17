package models

import (
	"time"
)

type NewBank struct {
	GlobalIdentifier string
	Name             string
	Address          string
	JurisdictionId   string
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

type Jurisdiction struct {
	Id   string
	Name string
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

type ComplianceCheck struct {
	Id                string
	OriginatorBankId  string
	BeneficiaryBankId string
	SenderId          int
	ReceiverId        int
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
	Id                            int
	PolicyTypeId                  int
	Owner                         string
	TransactionTypeId             int
	PolicyEnforcingJurisdictionId string
	OriginatingJurisdictionId     string
	BeneficiaryJurisdictionId     string
	Parameters                    string
	IsPrivate                     bool
	Latest                        bool
}

type PolicyAndItsType struct {
	Policy     NewPolicy
	PolicyType NewPolicyType
}

type NewFullPolicyModel struct {
	TransactionType NewTransactionType
	Policy          NewPolicy
	PolicyType      NewPolicyType
}

type NewTransactionPolicy struct {
	TransactionId string
	PolicyId      int
	Status        int
}

type NewTransactionProof struct {
	Id            uint64
	TransactionId uint64
	Proof         string
}

type NonInteractiveComplianceCheckProofRequest struct {
	ComplianceCheckId   string  `json:"compliance_check_id"`
	PolicyId            string  `json:"policy_id"`
	ParticipantsList    [][]int `json:"participants_list"`
	PublicSanctionsList [][]int `json:"pub_sanctions_list"`
}

type NonInteractiveComplianceCheckProofResponse struct {
	ID                    string                              `json:"id"`
	SanctionedCheckInput  NonInteractiveSanctionedCheckInput  `json:"sanctioned_check_input"`
	SanctionedCheckOutput NonInteractiveSanctionedCheckOutput `json:"sanctioned_check_output"`
	Status                string                              `json:"status"`
}

type NonInteractiveSanctionedCheckInput struct {
	ComplianceCheckID string  `json:"compliance_check_id"`
	PolicyID          string  `json:"policy_id"`
	ParticipantsList  [][]int `json:"participants_list"`
	PubSanctionsList  [][]int `json:"pub_sanctions_list"`
}

type NonInteractiveSanctionedCheckOutput struct {
	ComplianceCheckID    string `json:"compliance_check_id"`
	PolicyID             string `json:"policy_id"`
	ParticipantsListHash []int  `json:"participants_list_hash"`
	PubSanctionsListHash []int  `json:"pub_sanctions_list_hash"`
	NotSanctioned        bool   `json:"not_sanctioned"`
	Proof                []int  `json:"groth16_bn254_proof"`
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
