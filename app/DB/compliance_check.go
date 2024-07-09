package DB

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/errlog"
	"database/sql"
	"errors"
)

// TODO: currently it works with the "transaction" table, modify it to work with the "compliance check" table when that change is made in the database

// AddComplianceCheck creates a new compliance check and returns its id.
// Compliance check id is actually its SHA-256 hash. AddComplianceCheck internally
// ties up the newly created compliance check with the applicable policies for it.
func (h *DBHandler) AddComplianceCheck(complianceCheck models.ComplianceCheck) (string, error) {
	returnErr := errors.New("unsuccessful creation of compliance check")

	var err error

	if complianceCheck.Id == "" {
		// generate a unique compliance check id as a SHA-256 hash of its values
		complianceCheck.Id, err = common.GenerateHash(complianceCheck)
		if err != nil {
			errlog.Println(err)
			return "", returnErr
		}
	}

	// query to instert a new compliance check
	query := `INSERT INTO Transaction (Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId)
				  VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`
	row := h.db.QueryRow(query,
		sql.Named("p1", complianceCheck.Id),
		sql.Named("p2", complianceCheck.OriginatorBankId),
		sql.Named("p3", complianceCheck.BeneficiaryBankId),
		sql.Named("p4", complianceCheck.SenderId),
		sql.Named("p5", complianceCheck.ReceiverId),
		sql.Named("p6", complianceCheck.Currency),
		sql.Named("p7", complianceCheck.Amount),
		sql.Named("p8", complianceCheck.TransactionTypeId),
		sql.Named("p9", complianceCheck.LoanId))
	if row.Err() != nil {
		errlog.Println(err)

		return "", returnErr
	}

	// get all applicable policies for compliance check
	polices, err := h.GetPolicies(complianceCheck.OriginatorBankId, complianceCheck.BeneficiaryBankId, complianceCheck.TransactionTypeId)
	if err != nil {
		errlog.Println(err)

		return "", returnErr
	}

	// loop through all applicable policies and tie up them with the compliance check
	for _, policy := range polices {
		query = `INSERT INTO TransactionPolicy (TransactionId, PolicyId, Status) VALUES (@p1, @p2, @p3)`
		_, err := h.db.Exec(query,
			sql.Named("p1", complianceCheck.Id),
			sql.Named("p2", policy.Id),
			sql.Named("p3", 0))
		if err != nil {
			errlog.Println(err)
			return "", returnErr
		}
	}

	return complianceCheckId, nil
}
