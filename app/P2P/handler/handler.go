package handler

import (
	"bisgo/app/P2P/core"
	"bisgo/app/P2P/messages"
	"bisgo/app/models"
	"bisgo/common"
	"encoding/json"
	"log"
	"strconv"
)

type P2PHandler struct {
	*core.Core
}

func CreateP2PHandler(core *core.Core) *P2PHandler {
	return &P2PHandler{core}
}

func (h *P2PHandler) CreateTransaction(messageID int, payload []byte) {
	channel, err := messages.LoadChannel(messageID)

	if err != nil {
		// handle error
	}

	defer messages.RemoveChannel(messageID)

	// handler logic

	channel <- nil // send data to the listener

	var messageData common.TransactionDTO
	if err := json.Unmarshal(payload, &messageData); err != nil {
		log.Println(err.Error())
		return
	}

	transactionType, err := strconv.Atoi(messageData.TransactionType)
	if err != nil {
		log.Println(err.Error())
		return
	}

	transaction := models.Transaction{
		OriginatorBank:  h.DB.GetBankIdByIdentifier(messageData.OriginatorBankGlobalIdentifier),
		BeneficiaryBank: h.DB.GetBankIdByIdentifier(messageData.BeneficiaryBankGlobalIdentifier),
		Sender:          h.DB.GetBankClientId(messageData.SenderName),
		Receiver:        h.DB.GetBankClientId(messageData.ReceiverName),
		Currency:        messageData.Currency,
		Amount:          int(messageData.Amount),
		TypeId:          transactionType,
		LoanId:          int(messageData.LoanID),
	}

	transactionID := h.DB.InsertTransaction(transaction)
	h.DB.UpdateTransactionState(transactionID, 1)
}
