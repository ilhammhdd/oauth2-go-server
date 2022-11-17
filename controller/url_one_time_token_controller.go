package controller

import (
	"database/sql"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type UrlOneTimeToken struct {
	dbo DBOperator
}

func (uott UrlOneTimeToken) SelectURLOneTimeToken(clientID string, oneTimeToken string, signature string) (*entity.URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#UrlOneTimeToken.SelectURLOneTimeToken", callTraceFileRenderRegisterUserController)

	row, err := uott.dbo.QueryRow("SELECT id FROM clients WHERE client_id = ?", clientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "clients", "client_id"); detailedErr != nil {
		return nil, detailedErr
	}

	var clientsID uint64
	if err = row.Scan(&clientsID); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "clients.id", "id")
	}

	var urlOneTimeToken entity.URLOneTimeToken
	row, err = uott.dbo.QueryRow("SELECT * FROM url_one_time_tokens WHERE one_time_token = ? AND signature = ? AND clients_id = ?", oneTimeToken, signature, clientsID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "url_one_time_tokens", "one_time_token", "signature", "clients_id"); detailedErr != nil {
		return &urlOneTimeToken, detailedErr
	}

	if err = row.Scan(&urlOneTimeToken.ID, &urlOneTimeToken.CreatedAt, &urlOneTimeToken.UpdatedAt, &urlOneTimeToken.SoftDeletedAt, &urlOneTimeToken.Pk, &urlOneTimeToken.Sk, &urlOneTimeToken.OneTimeToken, &urlOneTimeToken.Signature, &urlOneTimeToken.URL, &urlOneTimeToken.ClientsID); err != nil && err != sql.ErrNoRows {
		return &urlOneTimeToken, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "url_one_time_tokens.*", "urlOneTimeToken")
	}
	return &urlOneTimeToken, nil
}

func (uott UrlOneTimeToken) DeleteURLOneTimeToken(oneTimeToken string, signature string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#UrlOneTimeToken.DeleteURLOneTimeToken", callTraceFileRenderRegisterUserController)

	_, err := uott.dbo.Command("DELETE FROM url_one_time_tokens WHERE one_time_token = ? AND signature = ?", oneTimeToken, signature)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBDelete, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "url_one_time_tokens")
	}
	return nil
}
