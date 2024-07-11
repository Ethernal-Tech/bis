package DB

import (
	"bisgo/app/models"
	"bisgo/errlog"
	"database/sql"
	"errors"
	"log"
)

func (wrapper *DBHandler) GetPolicyId(code string, jurisdictionId string) int {
	query := `SELECT Id FROM [Policy] WHERE Code = @p1 and JurisdictionId = @p2`

	rows, err := wrapper.db.Query(query, sql.Named("p1", code), sql.Named("p2", jurisdictionId))
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

// GetPolicies returns all policies that beneficiary bank and/or its/the central bank imposes to originator bank for the given transaction type.
func (h *DBHandler) GetPolicies(originatorBankId string, beneficiaryBankId string, transactionTypeId int) ([]models.NewPolicy, error) {
	returnErr := errors.New("unsuccessful obtainance of policies")

	// query to obtain jurisdiction of the originator bank
	query := `SELECT JurisdictionId FROM Bank WHERE GlobalIdentifier = @p1`

	var originatorJurisdiction string
	err := h.db.QueryRow(query, sql.Named("p1", originatorBankId)).Scan(&originatorJurisdiction)
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	// query to obtain jurisdiction of the beneficiary bank
	query = `SELECT JurisdictionId FROM Bank WHERE GlobalIdentifier = @p1`

	var beneficiaryJurisdiction string
	err = h.db.QueryRow(query, sql.Named("p1", beneficiaryBankId)).Scan(&beneficiaryJurisdiction)
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	// query to obtain the central bank of the beneficiary bank jurisdiction
	query = `SELECT GlobalIdentifier FROM Bank WHERE JurisdictionId = @p1 AND BankTypeId = 2`

	var beneficiaryCentralBankId string
	err = h.db.QueryRow(query, sql.Named("p1", beneficiaryJurisdiction)).Scan(&beneficiaryCentralBankId)
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	// query to obtain all policies that beneficiary jurisdiction (commercial bank + central bank) imposes to originator jurisdiction for the given transaction type
	query = `SELECT Id, PolicyTypeId, TransactionTypeId, Owner, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, Parameters, IsPrivate, Latest 
				FROM Policy WHERE OriginatingJurisdictionId = @p1 
				AND PolicyEnforcingJurisdictionId = @p2 
				AND TransactionTypeId = @p3
				AND Owner IN (@p4, @p5) 
				AND Latest = 1`
	rows, err := h.db.Query(query,
		sql.Named("p1", originatorJurisdiction),
		sql.Named("p2", beneficiaryJurisdiction),
		sql.Named("p3", transactionTypeId),
		sql.Named("p4", beneficiaryBankId),
		sql.Named("p5", beneficiaryCentralBankId))
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	defer rows.Close()

	policies := []models.NewPolicy{}

	// loop through the policies and append them to return policy list (slice)
	for rows.Next() {
		var policy models.NewPolicy
		err = rows.Scan(&policy.Id,
			&policy.PolicyTypeId,
			&policy.TransactionTypeId,
			&policy.Owner,
			&policy.PolicyEnforcingJurisdictionId,
			&policy.OriginatingJurisdictionId,
			&policy.Parameters,
			&policy.IsPrivate,
			&policy.Latest)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// GetPolicies returns all policies that apply to the passed jurisdiction for the given transaction type.
func (h *DBHandler) GetAppliedPolicies(jurisdictionId string, transactionTypeId int) ([]models.PolicyAndItsType, error) {
	returnErr := errors.New("unsuccessful obtainance of policies")

	// query to obtain all policies
	query := `SELECT p.Id, p.PolicyTypeId, p.Owner, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingJurisdictionId, p.OriginatingJurisdictionId, p.Parameters, p.IsPrivate, p.Latest
              FROM Policy p, PolicyType pt WHERE p.PolicyTypeId = pt.Id
              AND p.OriginatingJurisdictionId = @p1 AND p.TransactionTypeId = @p2 and p.Latest = 1`
	rows, err := h.db.Query(query,
		sql.Named("p1", jurisdictionId),
		sql.Named("p2", transactionTypeId))
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	defer rows.Close()

	policies := []models.PolicyAndItsType{}

	// loop through the policies and append them to return policy list (slice)
	for rows.Next() {
		var policy models.PolicyAndItsType
		err := rows.Scan(&policy.Policy.Id,
			&policy.Policy.PolicyTypeId,
			&policy.Policy.Owner,
			&policy.PolicyType.Code,
			&policy.PolicyType.Name,
			&policy.Policy.TransactionTypeId,
			&policy.Policy.PolicyEnforcingJurisdictionId,
			&policy.Policy.OriginatingJurisdictionId,
			&policy.Policy.Parameters,
			&policy.Policy.IsPrivate,
			&policy.Policy.Latest)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// GetPolicesByOwner provides the same functionality as [GetPolicies], but with additional filter on owner.
func (h *DBHandler) GetAppliedPoliciesByOwner(owner string, jurisdictionId string, transactionTypeId int) ([]models.PolicyAndItsType, error) {
	returnErr := errors.New("unsuccessful get of policies")

	// query to obtain all policies
	query := `SELECT p.Id, p.PolicyTypeId, p.Owner, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingJurisdictionId, p.OriginatingJurisdictionId, p.Parameters, p.IsPrivate, p.Latest
              FROM Policy p, PolicyType pt WHERE p.PolicyTypeId = pt.Id
              AND p.Owner = @p1 AND p.OriginatingJurisdictionId = @p2 AND p.TransactionTypeId = @p3 and p.Latest = 1`
	rows, err := h.db.Query(query,
		sql.Named("p1", owner),
		sql.Named("p2", jurisdictionId),
		sql.Named("p3", transactionTypeId))
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	defer rows.Close()

	policies := []models.PolicyAndItsType{}

	// loop through the policies and append them to return policy list (slice)
	for rows.Next() {
		var policy models.PolicyAndItsType
		err := rows.Scan(&policy.Policy.Id,
			&policy.Policy.PolicyTypeId,
			&policy.Policy.Owner,
			&policy.PolicyType.Code,
			&policy.PolicyType.Name,
			&policy.Policy.TransactionTypeId,
			&policy.Policy.PolicyEnforcingJurisdictionId,
			&policy.Policy.OriginatingJurisdictionId,
			&policy.Policy.Parameters,
			&policy.Policy.IsPrivate,
			&policy.Policy.Latest)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// GetPoliciesByComplianceCheckId returns all policies that apply to the given compliance check.
func (h *DBHandler) GetPoliciesByComplianceCheckId(complianceCheckId string) ([]models.PolicyAndItsType, error) {
	returnErr := errors.New("unsuccessful obtainance of policies")

	// query to obtain all policies
	query := `SELECT p.Id, p.PolicyTypeId, p.Owner, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingJurisdictionId, p.OriginatingJurisdictionId, p.Parameters, p.IsPrivate, p.Latest 
				FROM TransactionPolicy tp, Policy p, PolicyType pt 
				WHERE tp.PolicyId = p.Id AND p.PolicyTypeId = pt.Id AND tp.TransactionId = @p1`
	rows, err := h.db.Query(query, sql.Named("p1", complianceCheckId))
	if err != nil {
		errlog.Println(err)
		return nil, returnErr
	}

	defer rows.Close()

	policies := []models.PolicyAndItsType{}

	// loop through the policies and append them to return policy list (slice)
	for rows.Next() {
		var policy models.PolicyAndItsType
		err := rows.Scan(&policy.Policy.Id,
			&policy.Policy.PolicyTypeId,
			&policy.Policy.Owner,
			&policy.PolicyType.Code,
			&policy.PolicyType.Name,
			&policy.Policy.TransactionTypeId,
			&policy.Policy.PolicyEnforcingJurisdictionId,
			&policy.Policy.OriginatingJurisdictionId,
			&policy.Policy.Parameters,
			&policy.Policy.IsPrivate,
			&policy.Policy.Latest)
		if err != nil {
			errlog.Println(err)
			return nil, returnErr
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// UpdatePolicyStatus updates status of the policy for the given compliance check. Allowed values ​​for status are:
// 1. 0 - pending
// 2. 1 - passed
// 3. 2 - failed
func (h *DBHandler) UpdatePolicyStatus(complianceCheckId string, policyId int, status int) error {
	query := `UPDATE TransactionPolicy SET Status = @p1 WHERE TransactionId = @p2 AND PolicyId = @p3`

	_, err := h.db.Exec(query,
		sql.Named("p1", status),
		sql.Named("p2", complianceCheckId),
		sql.Named("p3", policyId))
	if err != nil {
		errlog.Println(err)
		return errors.New("unsuccessful update of policy status")
	}

	return nil
}

// GetPolicyStatus returns status of the policy for the given compliance check.
func (h *DBHandler) GetPolicyStatus(complianceCheckId string, policyId int) (int, error) {
	query := `SELECT Status FROM TransactionPolicy WHERE TransactionId = @p1 AND PolicyId = @p2`

	var status int
	err := h.db.QueryRow(query,
		sql.Named("p1", complianceCheckId),
		sql.Named("p2", policyId)).Scan(&status)
	if err != nil {
		errlog.Println(err)
		return -1, errors.New("unsuccessful obtainance of policy status")
	}

	return status, nil
}

// GetPolicyTypeByCode returns policy type with the given code.
func (h *DBHandler) GetPolicyTypeByCode(code string) (models.NewPolicyType, error) {
	query := `SELECT Id, Code, Name FROM PolicyType WHERE Code = @p1`

	var policyType models.NewPolicyType
	err := h.db.QueryRow(query,
		sql.Named("p1", code)).Scan(&policyType.Id,
		&policyType.Code,
		&policyType.Name)
	if err != nil {
		errlog.Println(err)
		return models.NewPolicyType{}, errors.New("unsuccessful obtainance of policy type")
	}

	return policyType, nil
}

// GetPolicyToProcessItsCheckResult is a special method to get the policy to process the result of its check.
// Due to its specificity, the given method is not intended to be used in any other case. Each policy can be
// uniquely identified by its id. However, unique identification can also be based on the following parameters:
// 1. policy type (id)
// 2. owner
// 3. transaction type (id)
// 4. originating jurisdiction (id)
// 5. isPrivate flag
func (h *DBHandler) GetPolicyToProcessItsCheckResult(policyTypeId int, owner string, transactionTypeId int, originatingJurisditionId string, isPrivate int) (models.NewPolicy, error) {
	query := `SELECT Id, PolicyTypeId, Owner, TransactionTypeId, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, Parameters, IsPrivate, Latest 
				FROM Policy WHERE PolicyTypeId = @p1 
				AND Owner = @p2 
				AND TransactionTypeId = @p3
				AND OriginatingJurisdictionId = @p4
				AND IsPrivate = @p5`

	var policy models.NewPolicy
	err := h.db.QueryRow(query,
		sql.Named("p1", policyTypeId),
		sql.Named("p2", owner),
		sql.Named("p3", transactionTypeId),
		sql.Named("p4", originatingJurisditionId),
		sql.Named("p5", isPrivate)).Scan(&policy.Id,
		&policy.PolicyTypeId,
		&policy.Owner,
		&policy.TransactionTypeId,
		&policy.PolicyEnforcingJurisdictionId,
		&policy.OriginatingJurisdictionId,
		&policy.Parameters,
		&policy.IsPrivate,
		&policy.Latest)
	if err != nil {
		errlog.Println(err)
		return models.NewPolicy{}, errors.New("unsuccessful obtainance of policy")
	}

	return policy, nil
}

func (wrapper *DBHandler) PoliciesFromJurisdiction(bankId string) []models.PolicyModel {
	query := `SELECT p.Id, pt.Code, pt.Name, tt.Name, p.Parameters 
					FROM Policy as p
					LEFT JOIN PolicyType as pt on p.PolicyTypeId = pt.Id
					LEFT JOIN TransactionType as tt on p.TransactionTypeId = tt.Id
					Where p.PolicyEnforcingJurisdictionId = (SELECT JurisdictionId FROM [Bank] Where GlobalIdentifier = @p1)`

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
	query := `SELECT p.Id, pt.Code, pt.Name, p.Parameters 
					FROM Policy as p
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

func (wrapper *DBHandler) GetPoliciesForTransaction(transactionID string) []models.PolicyAndItsType {
	query := `SELECT p.Id, p.PolicyTypeId, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingJurisdictionId, p.OriginatingJurisdictionId, p.Parameters, p.IsPrivate, p.Latest
              FROM Policy p
              JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
              JOIN TransactionPolicy tp ON p.Id = tp.PolicyId
              WHERE tp.TransactionId = @p1 AND p.Latest = 1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionID))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.PolicyAndItsType{}
	for rows.Next() {
		var policy models.PolicyAndItsType
		if err := rows.Scan(&policy.Policy.Id, &policy.Policy.PolicyTypeId, &policy.PolicyType.Code, &policy.PolicyType.Name, &policy.Policy.TransactionTypeId, &policy.Policy.PolicyEnforcingJurisdictionId, &policy.Policy.OriginatingJurisdictionId, &policy.Policy.Parameters, &policy.Policy.IsPrivate, &policy.Policy.Latest); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) GetAllPolicies() []models.NewFullPolicyModel {
	query := `SELECT p.Id, p.PolicyTypeId, pt.Code, pt.Name, p.TransactionTypeId, p.PolicyEnforcingJurisdictionId, p.OriginatingJurisdictionId, p.Parameters, p.IsPrivate, p.Latest, tt.Id, tt.Code, tt.Name
              FROM Policy p
              JOIN PolicyType pt ON p.PolicyTypeId = pt.Id
			  JOIN TransactionType tt ON p.TransactionTypeId = tt.Id
			  Where p.Latest = 1`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	policies := []models.NewFullPolicyModel{}
	for rows.Next() {
		var policy models.NewFullPolicyModel
		if err := rows.Scan(&policy.Policy.Id, &policy.Policy.PolicyTypeId, &policy.PolicyType.Code,
			&policy.PolicyType.Name, &policy.Policy.TransactionTypeId, &policy.Policy.PolicyEnforcingJurisdictionId,
			&policy.Policy.OriginatingJurisdictionId, &policy.Policy.Parameters, &policy.Policy.IsPrivate, &policy.Policy.Latest,
			&policy.TransactionType.Id, &policy.TransactionType.Code, &policy.TransactionType.Name); err != nil {
			log.Fatal(err)
		}

		policies = append(policies, policy)
	}
	return policies
}

func (wrapper *DBHandler) UpdateCFMPolicyAmount(amount uint64, policyId uint64) {
	query := `SELECT p.Id, 
					p.PolicyTypeId,
					p.TransactionTypeId,
					p.PolicyEnforcingJurisdictionId,
					p.OriginatingJurisdictionId,
					p.Parameters,
					p.IsPrivate,
					p.Latest 
				FROM Policy as p
				LEFT JOIN PolicyType as pt on p.PolicyTypeId = pt.Id
				Where p.Id = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", policyId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var policy models.NewPolicy
	for rows.Next() {
		if err := rows.Scan(&policy.Id, &policy.PolicyTypeId, &policy.TransactionTypeId, &policy.PolicyEnforcingJurisdictionId, &policy.OriginatingJurisdictionId,
			&policy.Parameters, &policy.IsPrivate, &policy.Latest); err != nil {
			log.Fatal(err)
		}
	}

	query = `UPDATE [Policy] Set Latest = 0 Where Id = @p1`

	_, err = wrapper.db.Exec(query,
		sql.Named("p1", policyId))
	if err != nil {
		log.Fatal(err)
	}

	insertQuery := `INSERT INTO [dbo].[Policy] (PolicyTypeId, TransactionTypeId, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, Parameters, IsPrivate, Latest) 
							VALUES (@p1, @p2, @p3, @p4, @p5, 0, 1)`
	_, err = wrapper.db.Exec(insertQuery,
		sql.Named("p1", policy.PolicyTypeId),
		sql.Named("p2", policy.TransactionTypeId),
		sql.Named("p3", policy.PolicyEnforcingJurisdictionId),
		sql.Named("p4", policy.OriginatingJurisdictionId),
		sql.Named("p5", amount))
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

// CreateOrGetPolicyType creates a new policy type and returns its id.
// If policy type already exists, method only returns id.
func (h *DBHandler) CreateOrGetPolicyType(code, name string) (int, error) {
	returnErr := errors.New("unsuccessful creation/obtainance of policy type/its id")

	// query to check if the policy type already exists in the system
	query := `SELECT Id FROM PolicyType WHERE Code = @p1 AND Name = @p2`

	var policyTypeId int
	err := h.db.QueryRow(query,
		sql.Named("p1", code),
		sql.Named("p2", name)).Scan(&policyTypeId)

	// if policy type (row) doesn't exist, error [sql.ErrNoRows] appears
	if err != nil {
		// check for potentially some other type of error
		if err != sql.ErrNoRows {
			errlog.Println(err)
			return -1, returnErr
		}

		// query to instert a new policy type and get its id
		query := `INSERT INTO PolicyType (Code, Name) OUTPUT INSERTED.Id VALUES (@p1, @p2)`
		err = h.db.QueryRow(query,
			sql.Named("p1", code),
			sql.Named("p2", name)).Scan(&policyTypeId)
		if err != nil {
			errlog.Println(err)
			return -1, returnErr
		}
	}

	return policyTypeId, nil
}

// CreatePolicy creates a new policy or updates the parameters of an existing one (read docs for [UpdatePolicyParameters]).
// Policy existence is checked against owner, originating jurisdiction, transaction type, policy type and isPrivate flag.
// Return values are as following:
// 1. id of the policy (created/updated one)
// 2. result status code (-1 - error; 0 - no action, 1 - created, 2 - updated)
// 3. error
func (h *DBHandler) CreateOrUpdatePolicy(policyTypeId int, owner string, transactionTypeId int, policyEnforcingJurisdictionId string, originatorJurisdictionId string, parameters string, isPrivate int) (int, int, error) {
	returnErr := errors.New("unsuccessful creation/update of policy")

	// query to check if the policy already exists in the system
	query := `SELECT p.Id FROM Policy p, PolicyType pt WHERE p.PolicyTypeId = pt.Id 
				AND p.PolicyTypeId = @p1
				AND p.Owner = @p2
				AND p.TransactionTypeId = @p3
				AND p.OriginatingJurisdictionId = @p4 
				AND p.IsPrivate = @p5
				AND p.Latest = 1`

	var policyId int
	err := h.db.QueryRow(query,
		sql.Named("p1", policyTypeId),
		sql.Named("p2", owner),
		sql.Named("p3", transactionTypeId),
		sql.Named("p4", originatorJurisdictionId),
		sql.Named("p5", isPrivate)).Scan(&policyId)

	// if policy (row) doesn't exist, error [sql.ErrNoRows] appears
	if err != nil {
		// check for potentially some other type of error
		if err != sql.ErrNoRows {
			errlog.Println(err)
			return -1, -1, returnErr
		}

		// query to instert a new policy and get its id
		query := `INSERT INTO Policy (PolicyTypeId, TransactionTypeId, Owner, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, Parameters, IsPrivate, Latest) 
					OUTPUT INSERTED.Id 
					VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, 1)`
		err = h.db.QueryRow(query,
			sql.Named("p1", policyTypeId),
			sql.Named("p2", transactionTypeId),
			sql.Named("p3", owner),
			sql.Named("p4", policyEnforcingJurisdictionId),
			sql.Named("p5", originatorJurisdictionId),
			sql.Named("p6", parameters),
			sql.Named("p7", isPrivate)).Scan(&policyId)
		if err != nil {
			errlog.Println(err)
			return -1, -1, returnErr
		}

		return policyId, 1, nil
	}

	// get the current parameters for the policy
	currentParameters, err := h.GetPolicyParameters(policyId)
	if err != nil {
		errlog.Println(err)
		return -1, -1, returnErr
	}

	// check for equality with the current parameters, if so, return
	if currentParameters == parameters {
		return policyId, 0, nil
	}

	// otherwise, update policy parameters
	policyId, err = h.UpdatePolicyParameters(policyId, parameters)
	if err != nil {
		errlog.Println(err)
		return -1, -1, returnErr
	}

	return policyId, 2, nil
}

// GetPolicyParameters returns the parameters of the selected (policyId) policy.
func (h *DBHandler) GetPolicyParameters(policyId int) (string, error) {
	query := `SELECT Parameters FROM Policy WHERE Id = @p1`

	var parameters string
	err := h.db.QueryRow(query,
		sql.Named("p1", policyId)).Scan(&parameters)
	if err != nil {
		errlog.Println(err)
		return "", errors.New("unsuccessful obtainance of policy parameters")
	}

	return parameters, nil
}

// UpdatePolicyParameters updates the parameters of the selected (policyId) policy.
// Update logic is not straightforward. Due to system specifics, the update actually takes the values ​​for
// all fields (expect the policyId and parameters) from the policy (row) selected by the passed (input) policyId and
// creates a new policy (row) with all these values ​​but a new policyId and passed parameters.
// Since this newly created policy is now marked as "latest", the old one has lost its status as latest.
// Error is returned if a policy selected by policyId doesn't exist or doesn't have latest tag set.
// Return value is the id of a newly created policy (row).
func (h *DBHandler) UpdatePolicyParameters(policyId int, parameters string) (int, error) {
	returnErr := errors.New("unsuccessful update of policy")

	// query to obtain the policy
	query := `SELECT Id, PolicyTypeId, TransactionTypeId, Owner, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, IsPrivate, Latest FROM Policy p WHERE p.Id = @p1`

	var policy models.NewPolicy
	err := h.db.QueryRow(query,
		sql.Named("p1", policyId)).Scan(&policy.Id,
		&policy.PolicyTypeId,
		&policy.TransactionTypeId,
		&policy.Owner,
		&policy.PolicyEnforcingJurisdictionId,
		&policy.OriginatingJurisdictionId,
		&policy.IsPrivate,
		&policy.Latest)

	// if policy doesn't exist, return an error
	if err != nil {
		errlog.Println(err)
		return -1, returnErr
	}

	// if policy selected by policyId isn't the latest one, return an error
	if !policy.Latest {
		return -1, errors.New("policy with a given policyId isn't the latest one")
	}

	// query to unset the "latest" flag of the previously latest policy
	query = `UPDATE Policy SET Latest = 0 Where Id = @p1`
	_, err = h.db.Exec(query, sql.Named("p1", policyId))
	if err != nil {
		errlog.Println(err)
		return -1, returnErr
	}

	// query to instert a new policy and get its id
	query = `INSERT INTO Policy (PolicyTypeId, TransactionTypeId, Owner, PolicyEnforcingJurisdictionId, OriginatingJurisdictionId, Parameters, IsPrivate, Latest) 
				OUTPUT INSERTED.Id 
				VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, 1)`
	err = h.db.QueryRow(query,
		sql.Named("p1", policy.PolicyTypeId),
		sql.Named("p2", policy.TransactionTypeId),
		sql.Named("p3", policy.Owner),
		sql.Named("p4", policy.PolicyEnforcingJurisdictionId),
		sql.Named("p5", policy.OriginatingJurisdictionId),
		sql.Named("p6", parameters),
		sql.Named("p7", policy.IsPrivate)).Scan(&policyId)
	if err != nil {
		errlog.Println(err)
		return -1, returnErr
	}

	return policyId, nil
}

func (wrapper *DBHandler) policyLatestStateChange(policyTypeId, transactionTypeId int, policyEnforcingJurisdictionId, originatorJurisdictionId, parameters string) error {
	// Check if Policy exists
	query := `SELECT Id FROM Policy WHERE PolicyTypeId = @p1 AND TransactionTypeId = @p2 AND PolicyEnforcingJurisdictionId = @p3
			  AND OriginatingJurisdictionId = @p4 AND Parameters <> @p5 AND Latest = 1`

	var policyIDs []uint

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", policyTypeId),
		sql.Named("p2", transactionTypeId),
		sql.Named("p3", policyEnforcingJurisdictionId),
		sql.Named("p4", originatorJurisdictionId),
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

func (wrapper *DBHandler) GetPolicyTypes() []models.NewPolicyType {
	query := `SELECT Id
					,Code
					,Name
				FROM PolicyType`

	policyTypes := []models.NewPolicyType{}
	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var policyType models.NewPolicyType
		if err := rows.Scan(&policyType.Id,
			&policyType.Code,
			&policyType.Name); err != nil {
			log.Fatal(err)
		}
		policyTypes = append(policyTypes, policyType)
	}
	return policyTypes
}
