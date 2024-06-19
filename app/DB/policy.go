package DB

import (
	"bisgo/app/models"
	"database/sql"
	"log"
)

func (wrapper *DBHandler) GetPolicyId(code string, countryId int) int {
	query := `SELECT Id FROM [Policy] WHERE Code = @p1 and CountryId = @p2`

	rows, err := wrapper.db.Query(query, sql.Named("p1", code), sql.Named("p2", countryId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var id int

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
	}

	return id
}

func (wrapper *DBHandler) GetPolicyById(policyID int) models.Policy {
	query := `SELECT Id, Code, Name FROM [Policy] WHERE Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", policyID))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var policy models.Policy

	for rows.Next() {
		if err := rows.Scan(&policy.Id, &policy.Code, &policy.Name); err != nil {
			log.Fatal(err)
		}
	}

	return policy
}

func (wrapper *DBHandler) GetPolices(bankId string, transactionTypeId int) []models.NewPolicy {
	query := `SELECT p.Id, p.PolicyTypeId, p.TransactionTypeId, p.PolicyEnforcingCountryId, p.OriginatingCountryId, p.Parameters
					FROM Policy p
					Join Country as c ON p.PolicyEnforcingCountryId = c.Id
					Where p.TransactionTypeId = @p2
						and p.PolicyEnforcingCountryId = (SELECT CountryId FROM [Bank] Where GlobalIdentifier = @p1)`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId),
		sql.Named("p2", transactionTypeId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.NewPolicy{}
	for rows.Next() {
		var policy models.NewPolicy
		if err := rows.Scan(&policy.Id, &policy.PolicyTypeId, &policy.TransactionTypeId, &policy.PolicyEnforcingCountryId, &policy.OriginatingCountryId, &policy.Parameters); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}

	return policies
}

func (wrapper *DBHandler) GetPolicesByCountryCode(originatingCountryCode string, transactionTypeId int) []models.NewPolicyModel {
	query := `SELECT p.Id, p.PolicyTypeId, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingCountryId, p.OriginatingCountryId, p.Parameters, p.IsPrivate, p.Latest
              FROM Policy p
              JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
              JOIN Country c ON p.OriginatingCountryId = c.Id
              WHERE c.Code = @p1 AND p.TransactionTypeId = @p2`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", originatingCountryCode),
		sql.Named("p2", transactionTypeId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.NewPolicyModel{}
	for rows.Next() {
		var policy models.NewPolicyModel
		if err := rows.Scan(&policy.Policy.Id, &policy.Policy.PolicyTypeId, &policy.PolicyType.Code, &policy.PolicyType.Name, &policy.Policy.TransactionTypeId, &policy.Policy.PolicyEnforcingCountryId, &policy.Policy.OriginatingCountryId, &policy.Policy.Parameters, &policy.Policy.IsPrivate, &policy.Policy.Latest); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) PoliciesFromCountry(bankId string) []models.PolicyModel {
	query := `SELECT p.Id, pt.Code, pt.Name, tt.Name, p.Parameters FROM Policy as p
					LEFT JOIN PolicyType as pt on p.PolicyTypeId = pt.Id
					LEFT JOIN TransactionType as tt on p.TransactionTypeId = tt.Id
					Where p.PolicyEnforcingCountryId = (SELECT CountryId FROM [Bank] Where GlobalIdentifier = @p1)`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.PolicyModel{}
	for rows.Next() {
		var policy models.PolicyModel
		if err := rows.Scan(&policy.Id, &policy.Code, &policy.Name, &policy.TransactionType, &policy.Parameter); err != nil {
			log.Fatal(err)
		}
		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) GetPolicy(policyId uint64) models.PolicyModel {
	query := `SELECT p.Id, pt.Code, pt.Name, p.Parameters FROM Policy as p
					LEFT JOIN PolicyType as pt on p.PolicyTypeId = pt.Id
					Where p.Id = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", policyId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var policy models.PolicyModel
	for rows.Next() {
		if err := rows.Scan(&policy.Id, &policy.Code, &policy.Name, &policy.Parameter); err != nil {
			log.Fatal(err)
		}
	}
	return policy
}

func (wrapper *DBHandler) GetPoliciesForTransaction(transactionID string) []models.NewPolicyModel {
	query := `SELECT p.Id, p.PolicyTypeId, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingCountryId, p.OriginatingCountryId, p.Parameters, p.IsPrivate, p.Latest
              FROM Policy p
              JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
              JOIN TransactionPolicy tp ON p.Id = tp.PolicyId
              WHERE tp.TransactionId = @p1 AND p.Latest = 1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionID))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.NewPolicyModel{}
	for rows.Next() {
		var policy models.NewPolicyModel
		if err := rows.Scan(&policy.Policy.Id, &policy.Policy.PolicyTypeId, &policy.PolicyType.Code, &policy.PolicyType.Name, &policy.Policy.TransactionTypeId, &policy.Policy.PolicyEnforcingCountryId, &policy.Policy.OriginatingCountryId, &policy.Policy.Parameters, &policy.Policy.IsPrivate, &policy.Policy.Latest); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) GetAllPolicies() []models.NewFullPolicyModel {
	query := `SELECT p.Id, p.PolicyTypeId, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingCountryId, p.OriginatingCountryId, p.Parameters, p.IsPrivate, p.Latest, tt.Id, tt.Code, tt.Name
              FROM Policy p
              JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
			  JOIN TransactionType tt ON p.TransactionTypeId = tt.Id`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.NewFullPolicyModel{}
	for rows.Next() {
		var policy models.NewFullPolicyModel
		if err := rows.Scan(&policy.Policy.Id, &policy.Policy.PolicyTypeId, &policy.PolicyType.Code,
			&policy.PolicyType.Name, &policy.Policy.TransactionTypeId, &policy.Policy.PolicyEnforcingCountryId,
			&policy.Policy.OriginatingCountryId, &policy.Policy.Parameters, &policy.Policy.IsPrivate, &policy.Policy.Latest,
			&policy.TransactionType.Id, &policy.TransactionType.Code, &policy.TransactionType.Name); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) UpdatePolicyAmount(amount uint64, policyId uint64) {
	query := `UPDATE [Policy] Set Parameters = @p1 Where Id = @p2`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", amount),
		sql.Named("p2", policyId))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBHandler) UpdatePolicyChecklist(checklist string, policyId uint64) {
	query := `UPDATE [Policy] Set Parameters = @p1 Where Id = @p2`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", checklist),
		sql.Named("p2", policyId))
	if err != nil {
		log.Fatal(err)
	}
}

// GetOrCreatePolicyType method to get policy type Id by Code and Name or create a new policy type if not exists
func (wrapper *DBHandler) GetOrCreatePolicyType(code, name string) uint {
	// Check if PolicyType exists
	query := `SELECT Id FROM PolicyType WHERE Code = @p1 AND Name = @p2`

	var policyTypeID uint
	err := wrapper.db.QueryRow(query, sql.Named("p1", code), sql.Named("p2", name)).Scan(&policyTypeID)
	if err != nil {
		if err == sql.ErrNoRows {
			// PolicyType does not exist, insert new PolicyType
			insertQuery := `INSERT INTO [dbo].[PolicyType] (Code, Name) OUTPUT INSERTED.Id VALUES (@p1, @p2)`
			err = wrapper.db.QueryRow(insertQuery, sql.Named("p1", code), sql.Named("p2", name)).Scan(&policyTypeID)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Other errors
			log.Fatal(err)
		}
	}

	return policyTypeID
}

// GetOrCreatePolicy method to get policy Id by its fields
func (wrapper *DBHandler) GetOrCreatePolicy(policyTypeId, transactionTypeId, policyEnforcingCountryId, originatorCountryId int, parameters string) uint {
	// Check if Policy exists
	query := `SELECT Id FROM Policy WHERE PolicyTypeId = @p1 AND TransactionTypeId = @p2 AND PolicyEnforcingCountryId = @p3
			  AND OriginatingCountryId = @p4 AND Parameters = @p5 AND Latest = 1`

	var policyID uint
	err := wrapper.db.QueryRow(query,
		sql.Named("p1", policyTypeId),
		sql.Named("p2", transactionTypeId),
		sql.Named("p3", policyEnforcingCountryId),
		sql.Named("p4", originatorCountryId),
		sql.Named("p5", parameters)).Scan(&policyID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Policy does not exist, insert new Policy
			insertQuery := `INSERT INTO [dbo].[Policy] (PolicyTypeId, TransactionTypeId, PolicyEnforcingCountryId, OriginatingCountryId, Parameters, IsPrivate, Latest) OUTPUT INSERTED.Id 
							VALUES (@p1, @p2, @p3, @p4, @p5, 0, 1)`
			err = wrapper.db.QueryRow(insertQuery,
				sql.Named("p1", policyTypeId),
				sql.Named("p2", transactionTypeId),
				sql.Named("p3", policyEnforcingCountryId),
				sql.Named("p4", originatorCountryId),
				sql.Named("p5", parameters)).Scan(&policyID)
			if err != nil {
				log.Fatal(err)
			}

			// If the policy already exists with the different parameters set it state to not be latest
			err = wrapper.policyLatestStateChange(policyTypeId, transactionTypeId, policyEnforcingCountryId, originatorCountryId, parameters)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// Other errors
			log.Fatal(err)
		}
	}

	return policyID
}

func (wrapper *DBHandler) policyLatestStateChange(policyTypeId, transactionTypeId, policyEnforcingCountryId, originatorCountryId int, parameters string) error {
	// Check if Policy exists
	query := `SELECT Id FROM Policy WHERE PolicyTypeId = @p1 AND TransactionTypeId = @p2 AND PolicyEnforcingCountryId = @p3
			  AND OriginatingCountryId = @p4 AND Parameters <> @p5 AND Latest = 1`

	var policyIDs []uint

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", policyTypeId),
		sql.Named("p2", transactionTypeId),
		sql.Named("p3", policyEnforcingCountryId),
		sql.Named("p4", originatorCountryId),
		sql.Named("p5", parameters))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var policyID uint
		if err := rows.Scan(&policyID); err != nil {
			return err
		}

		policyIDs = append(policyIDs, policyID)
	}

	for _, policyID := range policyIDs {
		query := `UPDATE [Policy] Set Latest = 0 Where Id = @p1`

		_, err := wrapper.db.Exec(query, sql.Named("p1", policyID))
		if err != nil {
			return err
		}
	}

	return nil
}
