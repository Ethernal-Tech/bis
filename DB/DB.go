package DB

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type DBWrapper struct {
	db *sql.DB
}

func InitDb() *DBWrapper {
	// Windows authentication
	sqldb, err := sql.Open("sqlserver", "sqlserver://@localhost:1434?database=BIS&trusted_connection=yes")
	// sqldb, err := sql.Open("sqlserver", "server=localhost;user id=SA;password=asdQWE123;port=1434;database=BIS")
	// sqldb, err := sql.Open("sqlserver", "sqlserver://testUser:123123@localhost:1434?database=BIS")

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
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var user BankEmployeeModel
		rows.Scan(&user.Name, &user.Username, &user.Password, &user.BankId, &user.BankName)
		return &user
	}

	return nil
}

func convertTxStatusDBtoPR(transaction *TransactionModel) *TransactionModel {

	switch transaction.Status {
	case "TransactionCreated":
		transaction.Status = "CREATED"
	case "PoliciesApplied":
		transaction.Status = "POLICIES APPLIED"
	case "ProofRequested":
		transaction.Status = "PROOF REQUESTED"
	case "ProofReceived":
		transaction.Status = "PROOF RECEIVED"
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

func convertTxStatusPRtoDB(transaction *TransactionModel) *TransactionModel {

	switch transaction.Status {
	case "CREATED":
		transaction.Status = "TransactionCreated"
	case "POLICIES APPLIED":
		transaction.Status = "PoliciesApplied"
	case "PROOF REQUESTED":
		transaction.Status = "ProofRequested"
	case "PROOF RECEIVED":
		transaction.Status = "ProofReceived"
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
	case "ProofRequested":
		history.Name = "PROOF REQUESTED"
	case "ProofReceived":
		history.Name = "PROOF RECEIVED"
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

func convertHistoryStatusPRtoDB(history *StatusHistoryModel) *StatusHistoryModel {

	switch history.Name {
	case "CREATED":
		history.Name = "TransactionCreated"
	case "POLICIES APPLIED":
		history.Name = "PoliciesApplied"
	case "PROOF REQUESTED":
		history.Name = "ProofRequested"
	case "PROOF RECEIVED":
		history.Name = "ProofReceived"
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
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	transactions := []TransactionModel{}
	for rows.Next() {
		var trnx TransactionModel
		rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdedntifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.Status)
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

	var trnx TransactionModel
	for rows.Next() {
		rows.Scan(&trnx.Id, &trnx.OriginatorBank, &trnx.BeneficiaryBank, &trnx.SenderGlobalIdentifier, &trnx.ReceiverGlobalIdedntifier, &trnx.SenderName, &trnx.ReceiverName, &trnx.Currency, &trnx.Amount, &trnx.Type, &trnx.TypeId)
	}
	rows.Close()

	query = `SELECT s.Name
					,th.Date
				FROM TransactionHistory th
				JOIN Status as s ON s.Id = th.StatusId
				Where Transactionid = @p1 Order by StatusId`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var statusHistory StatusHistoryModel
		rows.Scan(&statusHistory.Name, &statusHistory.Date)
		statusHistory = *convertHistoryStatusDBtoPR(&statusHistory)
		statusHistory.DateString = statusHistory.Date.Format("2006-01-02 15:04:05")
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	trnx.Status = trnx.StatusHistory[len(trnx.StatusHistory)-1].Name
	rows.Close()

	query = `SELECT p.Name
				FROM TransactionTypePolicy ttp
				JOIN Policy as p ON ttp.PolicyId = p.Id
				Where ttp.TransactionTypeId = (SELECT t.TransactionTypeId FROM [Transaction] as t
												Where t.Id = @p1)
					and ttp.CountryId = (SELECT b.CountryId FROM [Transaction] as t
										JOIN Bank as b ON b.Id = t.BeneficiaryBank
										Where t.Id = @p1)`

	rows, err = wrapper.db.Query(query, sql.Named("p1", transactionId))
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var policyName string
		rows.Scan(&policyName)
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

	var bankId uint64
	for rows.Next() {
		rows.Scan(&bankId)
	}

	return bankId
}

func (wrapper *DBWrapper) GetTransactionTypeId(transactionType string) int {
	query := `SELECT Id FROM [TransactionType] WHERE Name = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionType))

	if err != nil {
		log.Fatal(err)
	}

	var transactionTypeId int
	for rows.Next() {
		rows.Scan(&transactionTypeId)
	}

	return transactionTypeId
}

func (wrapper *DBWrapper) InsertTransaction(t Transaction) {
	query := `INSERT INTO [dbo].[Transaction] VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", t.OriginatorBank),
		sql.Named("p2", t.BeneficiaryBank),
		sql.Named("p3", t.Sender),
		sql.Named("p4", t.Receiver),
		sql.Named("p5", t.Currency),
		sql.Named("p6", t.Amount),
		sql.Named("p7", t.TypeId))
	if err != nil {
		log.Fatal(err)
	}
}

func (wrapper *DBWrapper) InsertTransactionPolicy(transactionId uint64, policies []int) {
	query := `INSERT INTO [dbo].[TransactionPolicy] VALUES (@p1, @p2)`

	for _, policy := range policies {
		_, err := wrapper.db.Exec(query,
			sql.Named("p1", transactionId),
			sql.Named("p2", policy))
		if err != nil {
			log.Fatal(err)
		}
	}
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

func (wrapper *DBWrapper) GetPolices(bankId uint64, transactionTypeId int) []PolicyModel {
	query := `SELECT c.Name, ttp.Amount, p.Name
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

	policies := []PolicyModel{}
	for rows.Next() {
		var policy PolicyModel
		rows.Scan(&policy.Country, &policy.Amount, &policy.Name)
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
