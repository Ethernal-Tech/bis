package DB

import (
	"bisgo/config"
	"bisgo/errlog"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// DBHandler represents a wrapper around the standard library's sql.DB type.
// It can be considered as a database connection pool.
type DBHandler struct {
	db *sql.DB
}

var handler DBHandler

func init() {
	var (
		db  *sql.DB
		err error
	)

	switch runtime.GOOS {
	case "linux":
		db, err = sql.Open("sqlserver", fmt.Sprintf("server=%s;user id=SA;password=%s;port=%s;database=%s", config.ResolveDBAddress(), config.ResolveDBPassword(), config.ResolveDBPort(), config.ResolveDBName()))
	case "windows":
		// windows authentication
		// db, err = sql.Open("sqlserver", fmt.Sprintf("sqlserver://@localhost:1434?database=%s&trusted_connection=yes", config.ResolveDBName()))

		// user authentication
		db, err = sql.Open("sqlserver", fmt.Sprintf("server=localhost;port=1433;database=%s;user id=testLogin;password=password", config.ResolveDBName()))
	default:
		errlog.Println(fmt.Errorf("\033[31m"+"DB handler is currently not supported for %s operating system."+"\033[31m", runtime.GOOS))
	}

	if err != nil {
		errlog.Println(err)
	}

	for i := 0; i < 60; i++ {
		err = db.Ping()
		if err != nil {
			errlog.Println(err)
		} else {
			break
		}

		time.Sleep(time.Second)
	}

	handler = DBHandler{db}
}

// GetDBHandler function return DBHandler that can be considered as a database connection pool.
// Currently supported operating systems:
//
// 1) linux
// 2) windows
func GetDBHandler() *DBHandler {
	return &handler
}

func Close() {
	handler.db.Close()
}
