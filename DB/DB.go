package DB

import (
	"database/sql"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DBWrapper struct {
	db *sql.DB
}

func InitDb() *DBWrapper {
	var (
		sqldb *sql.DB
		err   error
	)

	if runtime.GOOS == "linux" {
		sqldb, err = sql.Open("sqlserver", "server=localhost;user id=SA;password=Ethernal123;port=1433;database=BIS")
	} else {
		// Windows authentication
		sqldb, err = sql.Open("sqlserver", "sqlserver://@localhost:1434?database=BIS&trusted_connection=yes")
	}

	if err != nil {
		log.Panic(err)
	}

	err = sqldb.Ping()
	if err != nil {
		log.Panic(err)
	}

	wrapper := &DBWrapper{}
	wrapper.db = sqldb

	return wrapper
}

func (wrapper *DBWrapper) Login(username string, password string) *BankEmployeeModel {
	query := `SELECT [BankEmployee].Name Name
					,Username
					,Password
					,BankId
					,[Bank].Name BankName
				FROM [dbo].[BankEmployee], [dbo].[Bank]
				WHERE BankId = [Bank].Id AND Username = @p1 AND Password = @p2`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", username),
		sql.Named("p2", password))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var user BankEmployeeModel
		if err := rows.Scan(&user.Name, &user.Username, &user.Password, &user.BankId, &user.BankName); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		return &user
	}

	return nil
}

func (wrapper *DBWrapper) IsCentralBankEmploye(username string) bool {
	query := `SELECT CASE
					WHEN EXISTS (
						SELECT 1
						FROM BankEmployee AS be
						LEFT JOIN Bank AS b ON b.Id = be.BankId
						WHERE be.username = @p1 AND b.BankTypeId = @p2
					)
					THEN 'true'
					ELSE 'false'
				END AS CentralBankEmploye;`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", username),
		sql.Named("p2", CentralBank))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var CentralBankEmploye bool
		rows.Scan(&CentralBankEmploye)
		return CentralBankEmploye
	}

	return false
}

func convertTxStatusDBtoPR(transaction *TransactionModel) *TransactionModel {

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

// TODO: Remove if not neccessary
//
//nolint:unused
func convertTxStatusPRtoDB(transaction *TransactionModel) *TransactionModel {

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

func convertHistoryStatusDBtoPR(history *StatusHistoryModel) *StatusHistoryModel {

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

// TODO: Remove if not neccessary
//
//nolint:unused
func convertHistoryStatusPRtoDB(history *StatusHistoryModel) *StatusHistoryModel {

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

func (wrapper *DBWrapper) GetTransactionsForAddress(address uint64) []TransactionModel {
	query := `SELECT t.Id
					,ob.Name
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,bcs.Name
					,bcr.Name
					,t.Currency
					,t.Amount
					,s.Name
				FROM [Transaction] as t
				LEFT JOIN (SELECT MAX(StatusId) AS StatusId, Transactionid FROM TransactionHistory GROUP BY Transactionid) as th ON th.Transactionid = t.Id 
				LEFT JOIN [Status] as s ON s.Id = th.StatusId
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				WHERE t.OriginatorBank = @p1 OR t.BeneficiaryBank = @p2`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", address),
		sql.Named("p2", address))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	transactions := []TransactionModel{}
	for rows.Next() {
		var trnx TransactionModel
		if err := rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdedntifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.Status); err != nil {
			log.Println("Error scanning row:", err)
			return nil
		}
		trnx = *convertTxStatusDBtoPR(&trnx)
		transactions = append(transactions, trnx)
	}
	return transactions
}

func (wrapper *DBWrapper) GetTransactionHistory(transactionId uint64) TransactionModel {
	query := `SELECT t.Id
					,ob.Name
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,bcs.Name
					,bcr.Name
					,t.Currency
					,t.Amount
					,t.LoanId
					,ty.Code
					,ty.Name
					,ty.Id
				FROM [Transaction] as t
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				JOIN [TransactionType] as ty ON ty.Id = t.TransactionTypeId
				WHERE t.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trnx TransactionModel
	if rows.Next() {
		if err := rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdedntifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.LoanId, &trnx.TypeCode, &trnx.Type, &trnx.TypeId); err != nil {
			log.Fatal(err)
		}
	}

	query = `SELECT s.Name
					,th.Date
				FROM TransactionHistory th
				JOIN Status as s ON s.Id = th.StatusId
				Where Transactionid = @p1 Order by StatusId`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var statusHistory StatusHistoryModel
		if err := rows.Scan(&statusHistory.Name, &statusHistory.Date); err != nil {
			log.Fatal(err)
		}
		statusHistory = *convertHistoryStatusDBtoPR(&statusHistory)
		statusHistory.DateString = statusHistory.Date.Format("2006-01-02 15:04:05")
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	trnx.Status = trnx.StatusHistory[len(trnx.StatusHistory)-1].Name

	query = `SELECT p.Name
				FROM TransactionTypePolicy ttp
				JOIN Policy as p ON ttp.PolicyId = p.Id
				Where ttp.TransactionTypeId = (SELECT t.TransactionTypeId FROM [Transaction] as t
												Where t.Id = @p1)
					and ttp.CountryId = (SELECT b.CountryId FROM [Transaction] as t
										JOIN Bank as b ON b.Id = t.BeneficiaryBank
										Where t.Id = @p1)`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var policyName string
		if err := rows.Scan(&policyName); err != nil {
			log.Fatal(err)
		}
		trnx.Policies = append(trnx.Policies, policyName)
	}

	return trnx
}

func (wrapper *DBWrapper) GetBankId(bankName string) uint64 {
	query := `SELECT Id FROM [Bank] WHERE name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankName))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bankId uint64
	for rows.Next() {
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
	}

	return bankId
}

func (wrapper *DBWrapper) GetTransactionTypeId(transactionType string) int {
	query := `SELECT Id FROM [TransactionType] WHERE Code = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionType))

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var transactionTypeId int
	for rows.Next() {
		if err := rows.Scan(&transactionTypeId); err != nil {
			log.Fatal(err)
		}
	}

	return transactionTypeId
}

func (wrapper *DBWrapper) GetBankClientId(bankClientName string) uint64 {
	query := `SELECT Id FROM [BankClient] WHERE Name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankClientName))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bankClientId uint64
	for rows.Next() {
		if err := rows.Scan(&bankClientId); err != nil {
			log.Fatal(err)
		}
	}

	return bankClientId
}

func (wrapper *DBWrapper) InsertTransaction(t Transaction) uint64 {
	query := `INSERT INTO [dbo].[Transaction] OUTPUT INSERTED.Id VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`

	var insertedID uint64

	err := wrapper.db.QueryRow(query,
		sql.Named("p1", t.OriginatorBank),
		sql.Named("p2", t.BeneficiaryBank),
		sql.Named("p3", t.Sender),
		sql.Named("p4", t.Receiver),
		sql.Named("p5", t.Currency),
		sql.Named("p6", t.Amount),
		sql.Named("p7", t.TypeId),
		sql.Named("p8", t.LoanId)).Scan(&insertedID)

	if err != nil {
		log.Fatal(err)
	}

	polices := wrapper.GetPolices(t.BeneficiaryBank, t.TypeId)

	for _, policy := range polices {
		query = `INSERT INTO [dbo].[TransactionPolicyStatus] VALUES (@p1, @p2, @p3)`
		_, err := wrapper.db.Exec(query,
			sql.Named("p1", insertedID),
			sql.Named("p2", policy.Id),
			sql.Named("p3", 0))

		if err != nil {
			log.Fatal(err)
		}
	}

	return insertedID
}

func (wrapper *DBWrapper) UpdateTransactionState(transactionId uint64, state int) {
	query := `INSERT INTO [dbo].[TransactionHistory] VALUES (@p1, @p2, @p3)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", state),
		sql.Named("p3", time.Now()))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) Close() {
	wrapper.db.Close()
}

func (wrapper *DBWrapper) GetTransactionTypes() []TransactionType {
	query := `SELECT Id, Code, Name From TransactionType`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	types := []TransactionType{}
	for rows.Next() {
		var tType TransactionType
		if err := rows.Scan(&tType.Id, &tType.Code, &tType.Name); err != nil {
			log.Fatal(err)
		}
		types = append(types, tType)
	}
	return types
}

func (wrapper *DBWrapper) GetBanks() []BankModel {
	query := `SELECT b.Id, b.Name, c.Name
					From Bank b
					JOIN Country as c ON c.Id = b.CountryId`

	rows, err := wrapper.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	banks := []BankModel{}
	for rows.Next() {
		var bank BankModel
		if err := rows.Scan(&bank.Id, &bank.Name, &bank.Country); err != nil {
			log.Fatal(err)
		}
		banks = append(banks, bank)
	}
	return banks
}

func (wrapper *DBWrapper) GetBank(bankId uint64) Bank {
	query := `SELECT b.Id, b.Name, b.CountryId
					From Bank b
					WHERE b.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", bankId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var bank Bank

	for rows.Next() {
		if err := rows.Scan(&bank.Id, &bank.Name, &bank.CountryId); err != nil {
			log.Fatal(err)
		}
	}

	return bank
}

func (wrapper *DBWrapper) GetCountry(countryId uint) Country {
	query := `SELECT c.Id, c.Name, c.CountryCode
					From Country c
					WHERE c.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", countryId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var country Country

	for rows.Next() {
		if err := rows.Scan(&country.Id, &country.Name, &country.CountryCode); err != nil {
			log.Fatal(err)
		}
	}

	return country
}

func (wrapper *DBWrapper) GetTransactionPolicyStatuses(transactionId uint64) []TransactionPolicyStatus {
	query := `SELECT TransactionId, PolicyId, Status FROM [TransactionPolicyStatus] WHERE TransactionId = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var statuses []TransactionPolicyStatus

	for rows.Next() {
		var status TransactionPolicyStatus
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

func (wrapper *DBWrapper) GetPolices(bankId uint64, transactionTypeId int) []PolicyModel {
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

	policies := []PolicyModel{}
	for rows.Next() {
		var policy PolicyModel
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

func (wrapper *DBWrapper) InsertTransactionProof(transactionId uint64, value string) {
	query := `INSERT INTO [dbo].[TransactionProof] VALUES (@p1, @p2)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", transactionId),
		sql.Named("p2", value))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) GetPolicyId(code string, countryId int) int {
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

func (wrapper *DBWrapper) CheckCFM(receiverId uint64, countryId int) int64 {
	query := `SELECT GlobalIdentifier FROM BankClient Where Id = @p1`

	rows, err := wrapper.db.Query(query,
		sql.Named("p1", receiverId))

	if err != nil {
		log.Fatal(err)
	}

	var globalIdentifier string
	for rows.Next() {
		if err := rows.Scan(&globalIdentifier); err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()
	query = `SELECT b.Id FROM BankClient as bc Join (SELECT * FROM Bank Where CountryId = @p2) as b ON b.Id = bc.BankId	Where bc.GlobalIdentifier = @p1`

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", globalIdentifier),
		sql.Named("p2", countryId))

	if err != nil {
		log.Fatal(err)
	}

	var bankIds []string
	for rows.Next() {
		var bankId uint64
		if err := rows.Scan(&bankId); err != nil {
			log.Fatal(err)
		}
		bankIds = append(bankIds, strconv.Itoa(int(bankId)))
	}
	rows.Close()

	c := strings.Join(bankIds, ",")

	query = `SELECT
			(SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where Receiver = @p1 and BeneficiaryBank IN (` + c + `) and TransactionTypeId IN (1))
			-
			((SELECT ISNULL(SUM(Amount), 0)
			FROM [Transaction] t
            JOIN (SELECT TransactionId, StatusId FROM [TransactionHistory] WHERE StatusId = 7) as th on th.TransactionId = t.Id
			Where Sender = @p1 and OriginatorBank IN (` + c + `) and TransactionTypeId IN (2)))
			as difference`

	rows, err = wrapper.db.Query(query,
		sql.Named("p1", receiverId))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var amount int64
	for rows.Next() {
		if err := rows.Scan(&amount); err != nil {
			log.Fatal(err)
		}
	}

	return amount
}

func (wrapper *DBWrapper) PoliciesFromCountry(bankId uint64) []PolicyModel {
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

	policies := []PolicyModel{}
	for rows.Next() {
		var policy PolicyModel
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

func (wrapper *DBWrapper) GetPolicy(bankCountry string, policyId uint64) PolicyModel {
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

	var policy PolicyModel
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

func (wrapper *DBWrapper) UpdatePolicyAmount(policyId uint64, amount uint64) {
	query := `UPDATE [TransactionTypePolicy] Set Amount = @p2 Where PolicyId = @p1`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", policyId),
		sql.Named("p2", amount))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) UpdatePolicyChecklist(policyId uint64, checklist string) {
	query := `UPDATE [TransactionTypePolicy] Set Checklist = @p2 Where PolicyId = @p1`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", policyId),
		sql.Named("p2", checklist))
	if err != nil {
		log.Fatal(err)
	}
}
