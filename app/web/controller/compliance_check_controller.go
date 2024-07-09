package controller

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/errlog"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// AddComplianceCheck handles a web GET/POST "/addcompliancecheck" request. For a GET, it responds with a view for a new
// compliance check creation. On the other hand, POST indicates confirmation and that a compliance check should be created.
// As part of this process, a new compliance check is also sent over the p2p network to the beneficiary bank.
func (h *ComplianceCheckController) AddComplianceCheck(w http.ResponseWriter, r *http.Request) {
	if h.SessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodGet {

		// TODO: add logic together with a new addcompliancecheck view
		// TODO: modify json names for marshal/unmarshal

	} else if r.Method == http.MethodPost {

		data := struct {
			OriginatorGlobalIdentifier      string `json:"senderLei"`
			OriginatorName                  string `json:"senderName"`
			BeneficiaryGlobalIdentifier     string `json:"beneficiaryLei"`
			BeneficiaryName                 string `json:"beneficiaryName"`
			PaymentTypeId                   string `json:"paymentType"`
			TransactionTypeId               string `json:"transactionType"`
			Currency                        string `json:"currency"`
			Amount                          string `json:"amount"`
			BeneficiaryBankGlobalIdentifier string `json:"beneficiaryBank"`
		}{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		originatorBankGlobalIdentifier, err := h.DB.GetBankIdByName(h.SessionManager.GetString(r.Context(), "bankName"))
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		amount, err := strconv.Atoi(strings.Replace(data.Amount, ",", "", -1))
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		transactionTypeId, err := strconv.Atoi(data.TransactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		// TODO: potentially handle loan

		originatorId, err := h.DB.CreateOrGetBankClient(data.OriginatorGlobalIdentifier, data.OriginatorName, "", originatorBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		beneficiaryId, err := h.DB.CreateOrGetBankClient(data.BeneficiaryGlobalIdentifier, data.BeneficiaryName, "", data.BeneficiaryBankGlobalIdentifier)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		complianceCheck := models.ComplianceCheck{
			OriginatorBankId:  originatorBankGlobalIdentifier,
			BeneficiaryBankId: data.BeneficiaryBankGlobalIdentifier,
			SenderId:          originatorId,
			ReceiverId:        beneficiaryId,
			Currency:          data.Currency,
			Amount:            amount,
			TransactionTypeId: transactionTypeId,
			LoanId:            0,
		}

		complianceCheckId, err := h.DB.AddComplianceCheck(complianceCheck)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		transactionType, err := h.DB.GetTransactionTypeById(transactionTypeId)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		complianceCheckDTO := common.ComplianceCheckDTO{
			ComplianceCheckId:               complianceCheckId,
			OriginatorGlobalIdentifier:      data.OriginatorGlobalIdentifier,
			OriginatorName:                  data.OriginatorName,
			BeneficiaryGlobalIdentifier:     data.BeneficiaryBankGlobalIdentifier,
			BeneficiaryName:                 data.BeneficiaryName,
			OriginatorBankGlobalIdentifier:  originatorBankGlobalIdentifier,
			BeneficiaryBankGlobalIdentifier: data.BeneficiaryBankGlobalIdentifier,
			PaymentType:                     "",
			TransactionType:                 transactionType.Code,
			Amount:                          amount,
			Currency:                        data.Currency,
			SwiftBICCode:                    "",
			LoanId:                          0,
		}

		_, err = h.P2PClient.Send(data.BeneficiaryBankGlobalIdentifier, "new-compliance-check", complianceCheckDTO, 0)
		if err != nil {
			errlog.Println(err)

			http.Error(w, "Internal Server Error", 500)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
