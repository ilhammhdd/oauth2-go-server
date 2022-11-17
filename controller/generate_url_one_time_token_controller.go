package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileGenerateURLOneTimeTokenController = "/controller/generate_url_one_time_token_controller.go"

type initiatedURL struct {
	InitiatedURL string `json:"initiated_url"`
}

type GenerateURLOneTimeToken struct {
	dbo DBOperator
}

func (guott *GenerateURLOneTimeToken) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var callTraceFunc = fmt.Sprintf("%s#*GenerateURLOneTimeToken.SeveHTTP", callTraceFileGenerateURLOneTimeTokenController)

	authzToken, statusCode, authzTokenExistsResponseBody := IsAuthzTokenExists("Bearer", r.Header.Get("Authorization"), callTraceFunc)
	if authzToken == "" {
		w.WriteHeader(statusCode)
		w.Write(authzTokenExistsResponseBody)
		return
	}

	rules := make(map[string]uint)
	rules["client_id"] = adapter.RegexRandomID

	upv := restkit.URLQueryValidation{RegexRules: rules, Values: r.URL.Query()}
	if regexErrMsgs, valid := upv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen); !valid && regexErrMsgs != nil {
		responseBodyTmpl := entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexErrMsgs}
		jsonBody, err := json.Marshal(responseBodyTmpl)
		if err != nil {
			errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "response body template"))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)
			return
		}
	}

	usecase := usecase.GenerateURLOneTimeToken{
		ClientID: r.URL.Query().Get("client_id"), AuthzToken: authzToken, DBO: guott, ErrDescGen: errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc),
	}
	urlOneTimeToken, detailedErr := usecase.Generate(r.Host, r.URL.Path, r.URL.Query(), r.URL.Fragment)
	if errorkit.IsNotNilThenLog(detailedErr) {
		if detailedErr.ErrDescConst == entity.FlowErrUnauthorizedBearerAuthzToken || detailedErr.ErrDescConst == entity.FlowErrBearerAuthzTokenExpired {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(adapter.MakeGenerateURLOneTimeTokenErrResponse(nil, "", []error{detailedErr}))
		} else if !detailedErr.Flow {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(adapter.MakeGenerateURLOneTimeTokenErrResponse(nil, "", []error{detailedErr}))
		}
	} else {
		responseBody := initiatedURL{InitiatedURL: urlOneTimeToken.URL}
		jsonBody, err := json.Marshal(responseBody)
		if err != nil {
			errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "initiated url response body"))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(nil)
		} else {
			w.Header().Set("Location", urlOneTimeToken.URL)
			w.WriteHeader(http.StatusCreated)
			w.Write(jsonBody)
		}
	}
}

func (guott *GenerateURLOneTimeToken) SelectClientsBy(clientID string) (*entity.Client, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*GenerateURLOneTimeToken.SelectClientsBy", callTraceFileGenerateURLOneTimeTokenController)

	row, err := guott.dbo.QueryRow("SELECT id,client_id,client_secret,client_secret_expired_at FROM clients WHERE client_id = ?", clientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "clients", "client_id"); detailedErr != nil {
		return nil, detailedErr
	}

	var client entity.Client
	if err := row.Scan(&client.ID, &client.ClientID, &client.ClientSecret, &client.ClientSecretExpiredAt); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "clients", "client")
	}

	return &client, nil
}

func (guott *GenerateURLOneTimeToken) SelectCountURLOneTimeToken(clientsID uint64, url string) (uint32, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*GenerateURLOneTimeToken).SelectCountURLOneTimeToken", callTraceFileGenerateURLOneTimeTokenController)

	row, err := guott.dbo.QueryRow("SELECT COUNT(uott.id) FROM url_one_time_tokens uott JOIN clients c ON uott.clients_id = c.id WHERE c.id = ? AND uott.url = ?", clientsID, url)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "COUNT(url_one_time_tokens.id)", "clients.id", "url_one_time_tokens.id"); detailedErr != nil {
		return 0, detailedErr
	}

	var count uint32
	if err = row.Scan(&count); err != nil && err != sql.ErrNoRows {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "count(url_one_time_tokens.id)", "count")
	}
	return count, nil
}

func (guott *GenerateURLOneTimeToken) InsertURLOneTimeToken(urlOneTimeTokenEntity *entity.URLOneTimeToken) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#*GenerateURLOneTimeToken.InsertURLOneTimeToken", callTraceFileGenerateURLOneTimeTokenController)

	querySb := strings.Builder{}
	querySb.WriteString("INSERT INTO url_one_time_tokens (pk, sk, one_time_token, signature, url, clients_id) VALUES ")
	querySb.Write(sqlkit.GeneratePlaceHolder(6))
	if _, err := guott.dbo.Command(querySb.String(), urlOneTimeTokenEntity.Pk, urlOneTimeTokenEntity.Sk, urlOneTimeTokenEntity.OneTimeToken, urlOneTimeTokenEntity.Signature, urlOneTimeTokenEntity.URL, urlOneTimeTokenEntity.ClientsID); err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "url_one_time_tokens")
	}
	return nil
}

func NewGenerateURLOneTimeToken(db *sql.DB) *GenerateURLOneTimeToken {
	return &GenerateURLOneTimeToken{dbo: &sqlkit.DBOperation{DB: db}}
}
