package DB

import "bisgo/models"

func convertTxStatusDBtoPR(transaction *models.TransactionModel) *models.TransactionModel {
	switch transaction.Status {
	case "TransactionCreated":
		transaction.Status = "CREATED"
	case "PoliciesApplied":
		transaction.Status = "POLICIES APPLIED"
	case "ComplianceProofRequested":
		transaction.Status = "COMPLIANCE PROOF REQUESTED"
	case "ComplianceCheckPassed":
		transaction.Status = "COMPLIANCE CHECK PASSED"
	case "ProofInvalid":
		transaction.Status = "PROOF INVALID"
	case "AssetSent":
		transaction.Status = "ASSET SENT"
	case "TransactionCompleted":
		transaction.Status = "COMPLETED"
	case "TransactionCanceled":
		transaction.Status = "CANCELED"
	}
	return transaction
}

//nolint:unused
func convertTxStatusPRtoDB(transaction *models.TransactionModel) *models.TransactionModel {
	switch transaction.Status {
	case "CREATED":
		transaction.Status = "TransactionCreated"
	case "POLICIES APPLIED":
		transaction.Status = "PoliciesApplied"
	case "COMPLIANCE PROOF REQUESTED":
		transaction.Status = "ComplianceProofRequested"
	case "COMPLIANCE CHECK PASSED":
		transaction.Status = "ComplianceCheckPassed"
	case "PROOF INVALID":
		transaction.Status = "ProofInvalid"
	case "ASSET SENT":
		transaction.Status = "AssetSent"
	case "COMPLETED":
		transaction.Status = "TransactionCompleted"
	case "CANCELED":
		transaction.Status = "TransactionCanceled"
	}
	return transaction
}

func convertHistoryStatusDBtoPR(history *models.StatusHistoryModel) *models.StatusHistoryModel {
	switch history.Name {
	case "TransactionCreated":
		history.Name = "CREATED"
	case "PoliciesApplied":
		history.Name = "POLICIES APPLIED"
	case "ComplianceProofRequested":
		history.Name = "COMPLIANCE PROOF REQUESTED"
	case "ComplianceCheckPassed":
		history.Name = "COMPLIANCE CHECK PASSED"
	case "ProofInvalid":
		history.Name = "PROOF INVALID"
	case "AssetSent":
		history.Name = "ASSET SENT"
	case "TransactionCompleted":
		history.Name = "COMPLETED"
	case "TransactionCanceled":
		history.Name = "CANCELED"
	}
	return history
}

//nolint:unused
func convertHistoryStatusPRtoDB(history *models.StatusHistoryModel) *models.StatusHistoryModel {
	switch history.Name {
	case "CREATED":
		history.Name = "TransactionCreated"
	case "POLICIES APPLIED":
		history.Name = "PoliciesApplied"
	case "COMPLIANCE PROOF REQUESTED":
		history.Name = "ComplianceProofRequested"
	case "COMPLIANCE CHECK PASSED":
		history.Name = "ComplianceCheckPassed"
	case "PROOF INVALID":
		history.Name = "ProofInvalid"
	case "ASSET SENT":
		history.Name = "AssetSent"
	case "COMPLETED":
		history.Name = "TransactionCompleted"
	case "CANCELED":
		history.Name = "TransactionCanceled"
	}
	return history
}
