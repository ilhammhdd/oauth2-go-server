package controller

import (
	"database/sql"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type DBOperator interface {
	Command(stmt string, args ...interface{}) (sql.Result, error)
	Query(stmt string, args ...interface{}) (*sql.Rows, error)
	QueryRow(stmt string, args ...interface{}) (*sql.Row, error)
	QueryRowsToMap(stmt string, args ...interface{}) (*[]*map[string]interface{}, error)
}

func handleSelectDBNoRowsErr(err error, callTraceFunc string, args ...string) *errorkit.DetailedError {
	if err != nil && err != sql.ErrNoRows {
		return errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), args...)
	} else if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), args...)
	} else {
		return nil
	}
}
