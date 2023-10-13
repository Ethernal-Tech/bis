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
	// sqldb, err := sql.Open("sqlserver", "sqlserver://@localhost:1434?database=BIS&trusted_connection=yes")
	// sqldb, err := sql.Open("sqlserver", "server=localhost;user id=SA;password=asdQWE123;port=1434;database=BIS")
	sqldb, err := sql.Open("sqlserver", "sqlserver://testUser:123123@localhost:1434?database=BIS")

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

func (wrapper *DBWrapper) Login(username string, password string) BankEmployee {
	query := `SELECT Name
					,Username
					,Password
					,BankId 
				FROM [dbo].[BankEmployee] 
				WHERE Username = @p1 AND Password = @p2`

	rows, err := wrapper.db.Query(query, 
		sql.Named("p1", username), 
		sql.Named("p2", password))
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var user BankEmployee
	for rows.Next() {
		rows.Scan(&user.Name, &user.Username, &user.Password, &user.BankId)
	}

	return user
}

func (wrapper *DBWrapper) GetTransactionsForAddress(address uint64) []TransactionModel {
	query := `SELECT t.Id
					,ob.Name
					,bb.Name
					,bcs.GlobalIdentifier
					,bcr.GlobalIdentifier
					,t.Curency
					,t.Amount
					,ty.Name
					,s.Name
				FROM [Transaction] as t
				JOIN (SELECT MAX(StatusId) AS StatusId, Transactionid FROM TransactionHistory GROUP BY Transactionid) as th ON th.Transactionid = t.Id 
				JOIN [Status] as s ON s.Id = th.StatusId
				JOIN [Type] as ty ON ty.Id = t.Id
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
		rows.Scan(&trnx.Id, &trnx.BeneficiaryBank, &trnx.OriginatorBank, &trnx.Sender, &trnx.Receiver, &trnx.Curency, &trnx.Amount, &trnx.Type, &trnx.Status)
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
					,t.Curency
					,t.Amount
					,ty.Name
				FROM [Transaction] as t
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				JOIN [Type] as ty ON ty.Id = t.Id
				WHERE t.Id = @p1`

	rows, err := wrapper.db.Query(query, sql.Named("p1", transactionId))

	if err != nil {
		log.Fatal(err)
	}

	var trnx TransactionModel
	for rows.Next() {
		rows.Scan(&trnx.Id, &trnx.BeneficiaryBank, &trnx.OriginatorBank, &trnx.Sender, &trnx.Receiver, &trnx.Curency, &trnx.Amount, &trnx.Type)
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
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	trnx.Status = trnx.StatusHistory[len(trnx.StatusHistory)-1].Name
	rows.Close()

	query = `SELECT p.Name
				FROM TransactionPolicy tp
				JOIN Policy as p ON tp.PolicyId = p.Id
				Where Transactionid = @p1`

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

func (wrapper *DBWrapper) InsertTransaction(t Transaction) {
	query := `INSERT INTO [dbo].[Transaction] VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`

	_, err := wrapper.db.Exec(query,
		sql.Named("p1", t.OriginatorBank),
		sql.Named("p2", t.BeneficiaryBank),
		sql.Named("p3", t.Sender),
		sql.Named("p4", t.Receiver),
		sql.Named("p5", t.Curency),
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
