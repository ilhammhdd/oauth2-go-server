package external

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

var MariaDB *sql.DB

func OpenDBConnection(sourceName, user, password, database string) *sql.DB {
	dbDataSource := fmt.Sprintf("%s:%s@%s/%s?parseTime=true", user, password, sourceName, database)
	initDB, err := sql.Open("mysql", dbDataSource)
	errorkit.ErrorHandled(err, entity.StackTraceSize)

	return initDB
}
