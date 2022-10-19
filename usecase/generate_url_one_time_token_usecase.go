package usecase

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileGenerateURLOneTimeToken = "/usecase/generate_url_one_time_token.go"

type SelectClientByClientIDResult struct {
	ID           uint64
	ClientID     string
	ClientSecret string
}

type GenerateURLOneTimeTokenDBOperator interface {
	SelectClientByClientID(string) (*SelectClientByClientIDResult, *errorkit.DetailedError)
	SelectCountURLOneTimeToken(clientsID uint64, url string) (uint32, *errorkit.DetailedError)
	InsertURLOneTimeToken(*entity.URLOneTimeToken) *errorkit.DetailedError
}

type GenerateURLOneTimeToken struct {
	ClientID    string
	AccessToken string
	DBO         GenerateURLOneTimeTokenDBOperator
	ErrDescGen  errorkit.ErrDescGenerator
}

func (guott *GenerateURLOneTimeToken) Generate(host string, path string, query url.Values, fragment string) (*entity.URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#GenerateURLOneTimeToken", callTraceFileGenerateURLOneTimeToken)

	selectClientByClientIDResult, detailedErr := guott.DBO.SelectClientByClientID(guott.ClientID)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	}

	storedClientSecretRaw, err := base64.RawURLEncoding.DecodeString(selectClientByClientIDResult.ClientSecret)
	if err != nil {
		return nil, guott.handleErrors(errClientSecretDecode, callTraceFunc, err)
	}

	storedAccessTokenArr := blake2b.Sum256(storedClientSecretRaw)
	var accessToken []byte = storedAccessTokenArr[:]
	if guott.AccessToken != base64.RawURLEncoding.EncodeToString(accessToken) {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrUnauthorizedBearerAccessToken, guott.ErrDescGen)
	}

	entityURLOneTimeToken, detailedErr := entity.NewURLOneTimeToken(selectClientByClientIDResult.ID, selectClientByClientIDResult.ClientID, storedAccessTokenArr[:], host, path, query, fragment, guott.ErrDescGen)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	}

	existingURLOneTimeToken, detailedErr := guott.DBO.SelectCountURLOneTimeToken(selectClientByClientIDResult.ID, entityURLOneTimeToken.URL)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	}

	if existingURLOneTimeToken > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrExistsBasedOn, guott.ErrDescGen, "url_one_time_tokens", "clients.id", "url_one_time_tokens.url")
	}

	detailedErr = guott.DBO.InsertURLOneTimeToken(entityURLOneTimeToken)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	}

	return entityURLOneTimeToken, nil
}

func (guott *GenerateURLOneTimeToken) handleErrors(errType uint, callTraceFunc string, err error) *errorkit.DetailedError {
	switch errType {
	case errSelectClientByClientID:
		if err == sql.ErrNoRows {
			detailedErr := errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, guott.ErrDescGen, "clients", "client_id")
			errorkit.IsNotNilThenLog(detailedErr)
			return detailedErr
		} else {
			detailedErr := errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, guott.ErrDescGen)
			errorkit.IsNotNilThenLog(detailedErr)
			return detailedErr
		}
	case errGenerateRandomUUIDv4:
		detailedErr := errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrGenerateRandomUUIDv4, guott.ErrDescGen, "url_one_time_token")
		errorkit.IsNotNilThenLog(detailedErr)
		return detailedErr
	case errSelectCountURLOneTimeToken:
		detailedErr := errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBSelect, guott.ErrDescGen, "url_one_time_token")
		errorkit.IsNotNilThenLog(detailedErr)
		return detailedErr
	default:
		return nil
	}
}
