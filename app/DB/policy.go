package DB

import (
	"bisgo/app/models"
	"database/sql"
	"log"
	"strconv"
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

func (wrapper *DBHandler) GetPolices(bankId uint64, transactionTypeId int) []models.PolicyModel {
	query := `SELECT p.Id, c.Name, p.CountryId, p.Code, p.Name, ttp.Amount, ttp.Checklist
					FROM TransactionTypePolicy ttp
					JOIN Policy as p ON ttp.PolicyId = p.Id
					Join Country as c ON ttp.CountryId = c.Id
					Where ttp.TransactionTypeId = @p2
						and ttp.CountryId = (SELECT CountryId FROM [Bank] Where Id = @p1)`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId),
		sql.Named("p2", transactionTypeId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.PolicyModel{}
	for rows.Next() {
		var policy models.PolicyModel
		if err := rows.Scan(&policy.Id, &policy.Country, &policy.CountryId, &policy.Code, &policy.Name, &policy.Amount, &policy.Checklist); err != nil {
			log.Fatal(err)
		}

		if len(policy.Checklist) > 0 {
			policy.Parameter = policy.Checklist
		} else {
			policy.Parameter = strconv.FormatUint(policy.Amount, 10)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) PoliciesFromCountry(bankId uint64) []models.PolicyModel {
	query := `SELECT p.Id, c.Name, p.Code, p.Name, ttp.Amount, ttp.Checklist, tt.Name
				FROM TransactionTypePolicy ttp
				JOIN Policy as p ON ttp.PolicyId = p.Id
				Join Country as c ON ttp.CountryId = c.Id
				Join TransactionType as tt ON tt.Id = ttp.TransactionTypeId
				Where ttp.CountryId = (SELECT CountryId FROM [Bank] Where Id = @p1)`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.PolicyModel{}
	for rows.Next() {
		var policy models.PolicyModel
		if err := rows.Scan(&policy.Id, &policy.Country, &policy.Code, &policy.Name, &policy.Amount, &policy.Checklist, &policy.TransactionType); err != nil {
			log.Fatal(err)
		}

		if len(policy.Checklist) > 0 {
			policy.Parameter = policy.Checklist
		} else {
			policy.Parameter = strconv.FormatUint(policy.Amount, 10)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) GetPolicy(bankCountry string, policyId uint64) models.PolicyModel {
	query := `SELECT p.Id, c.Name, p.CountryId, p.Code, p.Name, ttp.Amount, ttp.Checklist
					FROM TransactionTypePolicy ttp
					JOIN Policy as p ON ttp.PolicyId = p.Id
					Join Country as c ON ttp.CountryId = c.Id
					Where ttp.PolicyId = @p2
						and ttp.CountryId = (SELECT Id FROM [Country] Where Name = @p1)`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", bankCountry),
		sql.Named("p2", policyId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var policy models.PolicyModel
	for rows.Next() {
		if err := rows.Scan(&policy.Id, &policy.Country, &policy.CountryId, &policy.Code, &policy.Name, &policy.Amount, &policy.Checklist); err != nil {
			log.Fatal(err)
		}

		if len(policy.Checklist) > 0 {
			policy.Parameter = policy.Checklist
		} else {
			policy.Parameter = strconv.FormatUint(policy.Amount, 10)
		}
	}

	return policy
}

func (wrapper *DBHandler) UpdatePolicyAmount(policyId uint64, amount uint64) {
	query := `UPDATE [TransactionTypePolicy] Set Amount = @p2 Where PolicyId = @p1`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", policyId),
		sql.Named("p2", amount))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBHandler) UpdatePolicyChecklist(policyId uint64, checklist string) {
	query := `UPDATE [TransactionTypePolicy] Set Checklist = @p2 Where PolicyId = @p1`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", policyId),
		sql.Named("p2", checklist))
	if err != nil {
		log.Fatal(err)
	}
}
