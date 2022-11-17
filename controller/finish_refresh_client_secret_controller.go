package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileFinishRefreshClientSecret = "/controller/finish_refresh_client_secret_controller.go"

type FinishRefreshClientSecret struct {
	dbo DBOperator
}

func (frcs FinishRefreshClientSecret) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var callTraceFunc = fmt.Sprintf("%s#FinishRefreshClientSecret.ServeHTTP", callTraceFileFinishRefreshClientSecret)

	uqv := restkit.URLQueryValidation{RegexRules: map[string]uint{
		"init_client_id": regexkit.RegexUUIDV4,
		"client_id":      adapter.RegexRandomID,
	},
		Values: r.URL.Query()}
	regexNoMatchMsgs, ok := uqv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNoMatchMsgs, "", nil))
		return
	}

	hpv := restkit.HeaderParamValidation{RegexRules: map[string]uint{"Authorization": adapter.RegexBearerAuthzToken}, Header: r.Header}
	regexNoMatchMsgs, ok = hpv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNoMatchMsgs, "", nil))
		return
	}

	reqBody, detailedErr := adapter.ReadRequestBody[entity.FinishRefreshClientSecretRequest](r, "refresh client secret")
	if detailedErr != nil {
		if detailedErr.Flow {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}
	errMsgs, ok := frcs.validateReqBody(reqBody)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(errMsgs, "", nil))
		return
	}

	theUsecase := usecase.NewFinishRefreshClientSecret(adapter.DetailedErrDescGen, frcs, reqBody, r.URL.Query().Get("client_id"))
	clientSecretExpiredAt, detailedErr := theUsecase.ValidateAndUpdateClientSecret(strings.Split(r.Header.Get("Authorization"), "Bearer ")[1])
	if detailedErr != nil {
		if detailedErr.Flow {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(adapter.MakeFinishRefreshClientSecretResponseBody(nil, "", nil, r.URL.Query().Get("init_client_id"), clientSecretExpiredAt))

}

func (frcs FinishRefreshClientSecret) validateReqBody(reqBody *entity.FinishRefreshClientSecretRequest) (*map[string][]string, bool) {
	var errMsgs map[string][]string = make(map[string][]string)

	if reqBody.InitClientID == "" {
		errMsgs["init_client_id"] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	}
	if reqBody.Basepoint == "" {
		errMsgs["basepoint"] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	}
	if reqBody.ClientPk == "" {
		errMsgs["client_pk"] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	}
	if reqBody.ServerPk == "" {
		errMsgs["server_pk"] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	}
	if !regexkit.RegexpCompiled[regexkit.RegexDateTimeRFC3339].Match([]byte(reqBody.SessionExpiredAt.Format(time.RFC3339Nano))) {
		errMsgs["session_expired_at"] = []string{adapter.RegexErrDescGen(regexkit.RegexDateTimeRFC3339)}
	}

	return &errMsgs, len(errMsgs) == 0
}

func (frcs FinishRefreshClientSecret) SelectClientRegistrationsBy(initClientID string) (*entity.ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#FinishRefreshClientSecret.SelectClientRegistrationsBy", callTraceFileFinishRefreshClientSecret)

	row, err := frcs.dbo.QueryRow("SELECT * FROM client_registrations WHERE init_client_id = ?", initClientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "client_registrations"); detailedErr != nil {
		return nil, detailedErr
	}

	var cr entity.ClientRegistration
	if err := row.Scan(&cr.ID, &cr.CreatedAt, &cr.UpdatedAt, &cr.SoftDeletedAt, &cr.InitClientID, &cr.Basepoint, &cr.ServerSK, &cr.ServerPK, &cr.SessionExpiredAt); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen, "client_registrations", "cr")
	}

	return &cr, nil
}

func (frcs FinishRefreshClientSecret) UpdateClientSecretAndExpiredAt(clientID string, initClientID, clientSecret string, clientSecretExpiredAt time.Time) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#FinishRefreshClientSecret.UpdateClientSecretAndExpiredAt", callTraceFileFinishRefreshClientSecret)

	_, err := frcs.dbo.Command("UPDATE clients SET init_client_id=?,client_secret=?,client_secret_expired_at=? WHERE client_id=?", initClientID, clientSecret, clientSecretExpiredAt, clientID)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBUpdate, adapter.DetailedErrDescGen, "client_secret", "client_secret_expired_at", "clients")
	}

	return nil
}

func (frcs FinishRefreshClientSecret) DeleteClientRegistration(initClientID string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#FinishRefreshClientSecret.DeleteClientRegistration", callTraceFileFinishRefreshClientSecret)

	if _, err := frcs.dbo.Command("DELETE FROM client_registrations WHERE init_client_id = ?", initClientID); err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBDelete, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations")
	}
	return nil
}

func NewFinishRefreshClientSecret(db *sql.DB) FinishRefreshClientSecret {
	return FinishRefreshClientSecret{dbo: sqlkit.DBOperation{DB: db}}
}
