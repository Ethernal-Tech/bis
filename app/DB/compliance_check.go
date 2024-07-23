package DB

import (
	"bisgo/app/models"
	"bisgo/common"
	"bisgo/errlog"
	"database/sql"
	"errors"
	"time"
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
	query := `INSERT INTO [Transaction] (Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId)
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
		errlog.Println(row.Err())

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
		query = `INSERT INTO TransactionPolicy (TransactionId, PolicyId, Status, AdditionalParameters, Description) VALUES (@p1, @p2, @p3, '', '')`
		_, err := h.db.Exec(query,
			sql.Named("p1", complianceCheck.Id),
			sql.Named("p2", policy.Id),
			sql.Named("p3", 0))
		if err != nil {
			errlog.Println(err)
			return "", returnErr
		}
	}

	return complianceCheck.Id, nil
}

// GetComplianceCheckById returns the compliance check with the given id.
func (h *DBHandler) GetComplianceCheckById(id string) (models.ComplianceCheck, error) {
	query := `SELECT Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId
              FROM [Transaction]
              WHERE Id = @p1`

	var complianceCheck models.ComplianceCheck
	err := h.db.QueryRow(query,
		sql.Named("p1", id)).Scan(&complianceCheck.Id,
		&complianceCheck.OriginatorBankId,
		&complianceCheck.BeneficiaryBankId,
		&complianceCheck.SenderId,
		&complianceCheck.ReceiverId,
		&complianceCheck.Currency,
		&complianceCheck.Amount,
		&complianceCheck.TransactionTypeId,
		&complianceCheck.LoanId)
	if err != nil {
		errlog.Println(err)
		return models.ComplianceCheck{}, errors.New("unsuccessful obtainance of compliance check")
	}

	return complianceCheck, nil
}

// UpdateComplianceCheckStatus updates the status of compliance check with the given id.
func (h *DBHandler) UpdateComplianceCheckStatus(id string, status int) error {
	query := `INSERT INTO TransactionHistory (TransactionId, StatusId, Date) VALUES (@p1, @p2, @p3)`

	_, err := h.db.Exec(query,
		sql.Named("p1", id),
		sql.Named("p2", status),
		sql.Named("p3", time.Now()))
	if err != nil {
		errlog.Println(err)
		return errors.New("unsuccessful update of compliance check state")
	}

	return nil
}

// GetAllSuccessfulComplianceChecks returns all successful compliance checks that the given entity has with entities from the given jurisdiction.
// Method works in two modes:
//
// 1. mode = 1
//   - entityId - denotes the id of the beneficiary
//   - jurisdictionId - denotes the id of the originator jurisdiction
//
// 2. mode = 2
//   - entityId - denotes the id of the originator
//   - jurisdictionId - denotes the id of the beneficiary jurisdiction
func (h *DBHandler) GetAllSuccessfulComplianceChecks(entityId int, jurisdictionId string, mode int) ([]models.ComplianceCheck, error) {
	returnErr := errors.New("unsuccessful obtainance of compliance checks")

	var query string

	if mode == 1 {
		// query to get all compliance checks in mode 1
		query = `SELECT Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId 
					FROM [Transaction] t, Bank b WHERE t.OriginatorBankId = b.GlobalIdentifier
					AND JurisdictionId = @p1
					AND ReceiverId = @p2`
	} else if mode == 2 {
		// query to get all compliance checks in mode 2
		query = `SELECT Id, OriginatorBankId, BeneficiaryBankId, SenderId, ReceiverId, Currency, Amount, TransactionTypeId, LoanId 
					FROM [Transaction] t, Bank b WHERE t.BeneficiaryBankId = b.GlobalIdentifier
					AND JurisdictionId = @p1
					AND SenderId = @p2`
	} else {
		return nil, returnErr
	}

	rows, err := h.db.Query(query, sql.Named("p1", jurisdictionId), sql.Named("p2", entityId))
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	defer rows.Close()

	allComplianceChecks := []models.ComplianceCheck{}

	// loop through the compliance checks and append them to the compliance check list (slice)
	for rows.Next() {
		var complianceCheck models.ComplianceCheck
		err := rows.Scan(&complianceCheck.Id,
			&complianceCheck.OriginatorBankId,
			&complianceCheck.BeneficiaryBankId,
			&complianceCheck.SenderId,
			&complianceCheck.ReceiverId,
			&complianceCheck.Currency,
			&complianceCheck.Amount,
			&complianceCheck.TransactionTypeId,
			&complianceCheck.LoanId)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		allComplianceChecks = append(allComplianceChecks, complianceCheck)
	}

	// compliance checks for which all policy checks were successful
	successfulComplianceChecks := []models.ComplianceCheck{}

complianceCheck:
	// loop through all compliance checks
	for _, compliancomplianceCheck := range allComplianceChecks {
		// get all policies for the given compliance check
		policies, err := h.GetPoliciesByComplianceCheckId(compliancomplianceCheck.Id)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		// loop through the policies
		for _, policy := range policies {
			policyStatus, err := h.GetPolicyStatus(compliancomplianceCheck.Id, policy.Policy.Id)
			if err != nil {
				errlog.Println(err)
				return nil, returnErr
			}

			// if any policy has non-successful status, dismiss compliance check
			if policyStatus != 1 {
				continue complianceCheck
			}
		}

		successfulComplianceChecks = append(successfulComplianceChecks, compliancomplianceCheck)
	}

	return successfulComplianceChecks, nil
}
