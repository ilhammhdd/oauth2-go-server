package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileInitRefreshClientSecret = "/controller/init_refresh_client_secret_controller.go"

type InitRefreshClientSecret struct {
	DBO DBOperator
}

func (ircs InitRefreshClientSecret) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var callTraceFunc = fmt.Sprintf("%s#InitRefreshClientSecret.ServeHTTP", callTraceFileInitRefreshClientSecret)

	uqv := restkit.URLQueryValidation{
		RegexRules: map[string]uint{
			"client_id": adapter.RegexRandomID, "init_client_id": regexkit.RegexNotEmpty,
		},
		Values: r.URL.Query()}
	regexNotMatch, ok := uqv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNotMatch, "", nil))
		return
	}

	hpv := restkit.HeaderParamValidation{RegexRules: map[string]uint{"Authorization": adapter.RegexBearerAuthzToken}, Header: r.Header}
	regexNotMatch, ok = hpv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNotMatch, "", nil))
		return
	}

	theUsecase := usecase.NewInitRefreshClientSecret(adapter.DetailedErrDescGen, ircs)
	cr, detailedErr := theUsecase.Initiate(r.URL.Query().Get("init_client_id"), r.URL.Query().Get("client_id"), strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if errorkit.IsNotNilThenLog(detailedErr) {
		if detailedErr.Flow {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusBadGateway)
		}
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", detailedErr))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(adapter.MakeInitiateRegisterClientResponseBody(nil, nil, cr))
}

func (ircs InitRefreshClientSecret) SelectClientsBy(clientID string) (*entity.Client, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#InitRefreshClientSecret.SelectClientsBy", callTraceFileInitRefreshClientSecret)

	row, err := ircs.DBO.QueryRow("SELECT client_secret, client_secret_expired_at FROM clients WHERE client_id = ?", clientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "clients"); detailedErr != nil {
		return nil, detailedErr
	}

	var client entity.Client
	if err := row.Scan(&client.ClientSecret, &client.ClientSecretExpiredAt); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen, "clients.client_secret, clients.client_secret_expired_at", "clientSecretAndExpired")
	}

	return &client, nil
}

func (ircs InitRefreshClientSecret) InsertClientRegistrations(cr *entity.ClientRegistration) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#InitRefreshClientSecret.InsertClientRegistrations", callTraceFileInitRefreshClientSecret)

	cr.SetSessionExpiredAt()

	_, err := ircs.DBO.Command(fmt.Sprintf("INSERT INTO client_registrations (init_client_id, basepoint, server_sk, server_pk, session_expired_at) VALUES %s", sqlkit.GeneratePlaceHolder(5)), cr.InitClientID, cr.Basepoint, cr.ServerSK, cr.ServerPK, cr.SessionExpiredAt)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, adapter.DetailedErrDescGen, "client_registrations")
	}
	return nil
}

func NewInitRefreshClientSecret(db *sql.DB) InitRefreshClientSecret {
	return InitRefreshClientSecret{sqlkit.DBOperation{DB: db}}
}
