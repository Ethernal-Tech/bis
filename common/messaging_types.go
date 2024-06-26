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

type PolicyRequestDTO struct {
	RequestID            string `json:"req_id"`
	BankGlobalIdentifier string `json:"bank_id"`
	TransactionType      string `json:"tx_type"`
}

type PolicyDTO struct {
	RequestID string `json:"req_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Params    string `json:"params"`
}

type SCLInitiationDTO struct {
	TransactionID string `json:"tx_id"`
	VMAddress     string `json:"vm_address"`
	PolicyID      uint   `json:"policy_id"`
}
