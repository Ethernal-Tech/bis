package common

type Peer struct {
	Name   string `json:"name"`
	PeerID string `json:"peer_id"`
}

type PassThruRequest struct {
	PeerID  string `json:"peer_id"`
	URI     string `json:"uri"`
	Payload []byte `json:"payload"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type TransactionDTO struct {
	TransactionID                   string `json:"tx_id"`
	SenderLei                       string `json:"sender_lei"`
	SenderName                      string `json:"sender_name"`
	ReceiverLei                     string `json:"receiver_lei"`
	ReceiverName                    string `json:"receiver_name"`
	OriginatorBankGlobalIdentifier  string `json:"orig_bank"`
	BeneficiaryBankGlobalIdentifier string `json:"benef_bank"`
	PaymentType                     string `json:"payment_type"`
	TransactionType                 string `json:"tx_type"`
	Amount                          uint64 `json:"amount"`
	Currency                        string `json:"curr"`
	SwiftBICCode                    string `json:"code"`
	LoanID                          uint64 `json:"load_id"`
}

type ComplianceCheckDTO struct {
	ComplianceCheckId               string `json:"compliance_check_id"`
	OriginatorGlobalIdentifier      string `json:"originator_lei"`
	OriginatorName                  string `json:"originator_name"`
	BeneficiaryGlobalIdentifier     string `json:"beneficiary_lei"`
	BeneficiaryName                 string `json:"beneficiary_name"`
	OriginatorBankGlobalIdentifier  string `json:"originator_bank"`
	BeneficiaryBankGlobalIdentifier string `json:"beneficiary_bank"`
	PaymentType                     string `json:"payment_type"`
	TransactionType                 string `json:"tx_type"`
	Amount                          int    `json:"amount"`
	Currency                        string `json:"currency"`
	SwiftBICCode                    string `json:"swift_code"`
	LoanId                          int    `json:"loan_id"`
}

// TODO: describe properties
type PolicyRequestDTO struct {
	Jurisdiction              string `json:"jurisdiction"`
	TransactionType           string `json:"tx_type"`
	RequesterGlobalIdentifier string `json:"requester"`
}

type PolicyResponseDTO struct {
	Policies []PolicyDTO `json:"policies"`
}

type PolicyDTO struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Params string `json:"params"`
	Owner  string `json:"owner"`
}

type CheckConfirmedDTO struct {
	CheckID   string `json:"tx_id"`
	VMAddress string `json:"vm_address"`
}

type CFMCheckDTO struct {
	TransctionID string `json:"tx_id"`
	Result       int    `json:"result"`
}
