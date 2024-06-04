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

func (wrapper *DBHandler) GetPolices(bankId string, transactionTypeId int) []models.PolicyModel {
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
