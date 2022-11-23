package controller

import (
	"database/sql"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type DBOperator interface {
	TxCommand(txCmdStmtArgs []*sqlkit.TxCmdStmtArgs) ([]*sql.Result, error)
	Command(stmt string, args ...interface{}) (sql.Result, error)
	Query(stmt string, args ...interface{}) (*sql.Rows, error)
	QueryRow(stmt string, args ...interface{}) (*sql.Row, error)
	QueryRowsToMap(stmt string, args ...interface{}) (*[]*map[string]interface{}, error)
}

func handleSelectTableErr(err error, callTraceFunc string, errArgs ...string) *errorkit.DetailedError {
	if err != nil && err == sql.ErrNoRows {
		return errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), errArgs...)
	} else if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), errArgs...)
	} else {
		return nil
	}
}
