package DB

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func InitDb() *sql.DB {
	sqldb, err := sql.Open("sqlserver", "sqlserver://testUser:password@localhost:1434?database=BIS")
	if err != nil {
		log.Panic(err)
	}

	err = sqldb.Ping()
	if err != nil {
		log.Panic(err)
	}

	return sqldb
}
