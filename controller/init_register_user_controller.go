package controller

import (
	"database/sql"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/sqlkit"
)

type InitRegisterUserController struct {
	dbo DBOperator
}

func (iruc InitRegisterUserController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func NewInitRegisterUserController(db *sql.DB) InitRegisterUserController {
	return InitRegisterUserController{dbo: sqlkit.DBOperation{DB: db}}
}
