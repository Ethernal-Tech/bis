package DB

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type DBWrapper struct {
	db	*sql.DB
}

func InitDb() *DBWrapper {
	sqldb, err := sql.Open("sqlserver", "sqlserver://testUser:password@localhost:1434?database=BIS")
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

func (wrapper *DBWrapper) GetTransactionsForAddress(address uint64) {

}

func (wrapper *DBWrapper) GetTransactionHistory(transactionId uint64) {

}

func (wrapper *DBWrapper) InsertTransaction(t Transaction) {
	
}

func (wrapper *DBWrapper) InsertTransactionPolicy(transactionId uint64, policies []int) {
	
}

func (wrapper *DBWrapper) UpdateTransactionState(transactionId uint64, state TransactionStates) {
	
}
