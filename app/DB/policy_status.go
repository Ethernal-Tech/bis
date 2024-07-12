package DB

import (
	"bisgo/app/models"
	"database/sql"
	"log"
)

func (wrapper *DBHandler) GetTransactionPolicyStatuses(transactionId string) []models.NewTransactionPolicy {
	query := `SELECT TransactionId, PolicyId, Status FROM [TransactionPolicy] WHERE TransactionId = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []models.NewTransactionPolicy

	for rows.Next() {
		var status models.NewTransactionPolicy
		if err := rows.Scan(&status.TransactionId, &status.PolicyId, &status.Status); err != nil {
			log.Fatal(err)
		}
		statuses = append(statuses, status)
	}

	return statuses
}

func (wrapper *DBHandler) UpdateTransactionPolicyStatus(transactionId string, policyId int, status int) {
	query := `UPDATE [TransactionPolicy] Set Status = @p3 Where TransactionId = @p1 and PolicyId = @p2`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", policyId),
		sql.Named("p3", status))
	if err != nil {
		log.Fatal(err)
	}
}
