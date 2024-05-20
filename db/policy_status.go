package db

import (
	"bisgo/common"
	"database/sql"
	"log"
)

func (wrapper *DBWrapper) GetTransactionPolicyStatuses(transactionId uint64) []common.TransactionPolicyStatus {
	query := `SELECT TransactionId, PolicyId, Status FROM [TransactionPolicyStatus] WHERE TransactionId = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []common.TransactionPolicyStatus

	for rows.Next() {
		var status common.TransactionPolicyStatus
		if err := rows.Scan(&status.TransactionId, &status.PolicyId, &status.Status); err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)
	}

	return statuses
}

func (wrapper *DBWrapper) GetTransactionPolicyStatus(transactionId uint64, policyId int) int {
	query := `SELECT Status FROM [TransactionPolicyStatus] WHERE TransactionId = @p1 and PolicyId = @p2`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId), sql.Named("p2", policyId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var status int

	for rows.Next() {
		if err := rows.Scan(&status); err != nil {
			log.Fatal(err)
		}
	}

	return status
}

func (wrapper *DBWrapper) UpdateTransactionPolicyStatus(transactionId uint64, policyId int, status int) {
	query := `UPDATE [TransactionPolicyStatus] Set Status = @p3 Where TransactionId = @p1 and PolicyId = @p2`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", policyId),
		sql.Named("p3", status))
	if err != nil {
		log.Fatal(err)
	}
}
