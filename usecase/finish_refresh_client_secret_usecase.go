package usecase

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileFinishRefreshClientSecret = "/usecase/finish_refresh_client_secret_usecase.go"

type FinishRefreshClientSecretDBO interface {
	SelectClientRegistrationsBy(initClientID string) (*entity.ClientRegistration, *errorkit.DetailedError)
	UpdateClientSecretAndExpiredAt(clientID string, newInitClientID, clientSecret string, clientSecretExpiredAt time.Time) *errorkit.DetailedError
	DeleteClientRegistration(initClientID string) *errorkit.DetailedError
}

type FinishRefreshClientSecret struct {
	errDescGen  errorkit.ErrDescGenerator
	dbo         FinishRefreshClientSecretDBO
	RequestBody *entity.FinishRefreshClientSecretRequest
	ClientID    string
}

func (frcs FinishRefreshClientSecret) ValidateAndUpdateClientSecret(authzToken string) (time.Time, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#FinishRefreshClientSecret.CalculateClientSecret", callTraceFileFinishRefreshClientSecret)

	defer frcs.dbo.DeleteClientRegistration(frcs.RequestBody.InitClientID)

	cr, detailedErr := frcs.dbo.SelectClientRegistrationsBy(frcs.RequestBody.InitClientID)
	if detailedErr != nil {
		return time.Time{}, detailedErr
	}

	clientSecretRaw, detailedErr := cr.CalculateClientSecret(frcs.RequestBody.ClientPk, frcs.errDescGen)
	if detailedErr != nil {
		return time.Time{}, detailedErr
	}
	clientSecretChecksum := blake2b.Sum256(clientSecretRaw)
	clientSecret := base64.RawURLEncoding.EncodeToString(clientSecretRaw)

	if authzToken != base64.RawURLEncoding.EncodeToString(clientSecretChecksum[:]) {
		return time.Time{}, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrUnauthorizedBearerAuthzToken, frcs.errDescGen)
	}

	clientSecretExpiredAt, detailedErr := entity.DetermineClientSecretExpiredAt(frcs.errDescGen)
	if detailedErr != nil {
		return time.Time{}, detailedErr
	}

	detailedErr = frcs.dbo.UpdateClientSecretAndExpiredAt(frcs.ClientID, frcs.RequestBody.InitClientID, clientSecret, clientSecretExpiredAt)
	if detailedErr != nil {
		return time.Time{}, detailedErr
	}

	return clientSecretExpiredAt, nil
}

func NewFinishRefreshClientSecret(errDescGen errorkit.ErrDescGenerator, dbo FinishRefreshClientSecretDBO, requestBody *entity.FinishRefreshClientSecretRequest, clientID string) FinishRefreshClientSecret {
	return FinishRefreshClientSecret{errDescGen, dbo, requestBody, clientID}
}
