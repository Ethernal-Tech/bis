package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
	"log"
)

func (wrapper *DBHandler) GetTransactionPolicyStatuses(transactionId string) ([]models.NewTransactionPolicy, error) {
	returnErr := errors.New("unsuccessful query of transaction policy statuses")
	query := `SELECT TransactionId, PolicyId, Status, AdditionalParameters, Description FROM [TransactionPolicy] WHERE TransactionId = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []models.NewTransactionPolicy

	for rows.Next() {
		var status models.NewTransactionPolicy
		if err := rows.Scan(
			&status.TransactionId,
			&status.PolicyId,
			&status.Status,
			&status.AdditionalParameters,
			&status.Description); err != nil {
			errlog.Println(err)
			return statuses, returnErr
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (wrapper *DBHandler) UpdateTransactionPolicyStatus(transactionId string, policyId int, status int, description string) {
	query := `UPDATE [TransactionPolicy] SET Status = @p3, Description = @p4 WHERE TransactionId = @p1 AND PolicyId = @p2`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", policyId),
		sql.Named("p3", status),
		sql.Named("p4", description))
	if err != nil {
		log.Fatal(err)
	}
}
