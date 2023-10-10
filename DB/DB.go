package DB

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type DBWrapper struct {
	db *sql.DB
}

func InitDb() *DBWrapper {
	// Windows authentication
	//sqldb, err := sql.Open("sqlserver", "sqlserver://@localhost:1434?database=BIS&trusted_connection=yes")
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
				JOIN  (SELECT MAX(StatusId) AS StatusId, Transactionid FROM TransactionHistory GROUP BY Transactionid) as th ON th.Transactionid = t.Id 
				JOIN [Status] as s ON s.Id = th.StatusId
				JOIN [Type] as ty ON ty.Id = t.Id
				JOIN Bank as ob ON ob.Id = t.OriginatorBank
				JOIN Bank as bb ON bb.Id = t.BeneficiaryBank
				JOIN BankClient as bcs ON bcs.Id = t.Sender
				JOIN BankClient as bcr ON bcr.Id = t.Receiver
				WHERE t.OriginatorBank = $1 OR t.BeneficiaryBank = $2`

	rows, err := wrapper.db.Query(query, address, address)
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
				WHERE t.Id = $1`

	rows, err := wrapper.db.Query(query, transactionId)

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
				Where Transactionid = $1 Order by StatusId`

	rows, err = wrapper.db.Query(query, transactionId)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var statusHistory StatusHistoryModel
		rows.Scan(&statusHistory.Name, &statusHistory.Date)
		trnx.StatusHistory = append(trnx.StatusHistory, statusHistory)
	}
	return trnx
}

func (wrapper *DBWrapper) InsertTransaction(t Transaction) {

}

func (wrapper *DBWrapper) InsertTransactionPolicy(transactionId uint64, policies []int) {

}

func (wrapper *DBWrapper) UpdateTransactionState(transactionId uint64, state int) {

}

func (wrapper *DBWrapper) Close() {
	wrapper.db.Close()
}
