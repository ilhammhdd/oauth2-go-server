package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileInitRegisterClientController = "/controller/init_register_client_controller.go"

type InitClientRegistration struct {
	dbo DBOperator
}

func (icr InitClientRegistration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var callTraceFunc = fmt.Sprintf("%s#InitClientRegistration.ServeHTTP", callTraceFileInitRegisterClientController)
	var errs []error

	rules := make(map[string]uint)
	rules["register_type"] = adapter.RegexClientRegisterType
	rules["init_client_id_checksum"] = regexkit.RegexNotEmpty

	upv := &restkit.URLParamValidation{RegexRules: rules, Values: r.URL.Query()}

	w.Header().Set("Content-Type", "application/json")

	if regexErrMsgs, valid := upv.Validate(errorkit.ErrDescGeneratorFunc(adapter.GenerateRegexErrDesc)); !valid && regexErrMsgs != nil {
		response, err := json.Marshal(entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexErrMsgs})
		if err != nil {
			errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "response body template"))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	cr, err := usecase.InitiateClientRegistration(r.URL.Query().Get("init_client_id_checksum"), icr, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc))
	if errorkit.IsNotNilThenLog(err) {
		errs = append(errs, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(adapter.MakeRegisterInitiateClientResponseBody(nil, errs, cr))
}

func (icr InitClientRegistration) InsertIgnoreDBO(cr *entity.ClientRegistration) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#InitClientRegistration.InsertIgnoreDBO", callTraceFileInitRegisterClientController)
	cr.SetSessionExpiredAt()
	_, err := icr.dbo.Command(fmt.Sprintf("INSERT IGNORE INTO client_registrations(init_client_id_checksum,basepoint,server_sk,server_pk,session_expired_at) VALUES %s", string(sqlkit.GeneratePlaceHolder(5))), cr.InitClientIDChecksum, cr.Basepoint, cr.ServerSK, cr.ServerPK, cr.SessionExpiredAt)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations")
	}
	return nil
}

func (icr InitClientRegistration) SelectCountBy(initClientIdChecksum string) (int, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#InitClientRegistration.SelectCountBy", callTraceFileInitRegisterClientController)
	row, err := icr.dbo.QueryRow("SELECT COUNT(id) FROM client_registrations WHERE init_client_id_checksum = ?", initClientIdChecksum)
	if err != nil && err != sql.ErrNoRows {
		return -1, errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations", "COUNT(client_registrations.id)")
	} else if err != nil {
		return -1, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations")
	}

	var counted int = -1
	err = row.Scan(&counted)
	if err != nil {
		return -1, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "COUNT(client_registrations.id)", "counted")
	}
	return counted, nil
}

func (icr InitClientRegistration) SelectCountClientsBy(initClientIDChecksum string) (int, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#InitClientRegistration.SelectCountClientsBy", callTraceFileInitRegisterClientController)
	row, err := icr.dbo.QueryRow("SELECT COUNT(id) FROM clients WHERE init_client_id_checksum = ?", initClientIDChecksum)
	if err != nil && err != sql.ErrNoRows {
		return -1, errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "COUNT(clients.id)", "init_client_id_checksum")
	} else if err != nil {
		return -1, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "clients")
	}

	var counted int
	if err := row.Scan(&counted); err != nil {
		return -1, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "count(clients.id)", "counted")
	}
	return counted, nil
}

func NewInitiateClientRegistration(db *sql.DB) InitClientRegistration {
	return InitClientRegistration{dbo: &sqlkit.DBOperation{DB: db}}
}
