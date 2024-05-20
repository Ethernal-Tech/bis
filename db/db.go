package db

import (
	"database/sql"
	"log"
	"runtime"

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

	wrapper := &DBWrapper{
		db: sqldb,
	}
	wrapper.db = sqldb

	return wrapper
}

func (wrapper *DBWrapper) Close() error {
	return wrapper.db.Close()
}
