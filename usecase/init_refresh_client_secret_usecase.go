package usecase

import (
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileInitRefreshClientSecret = "/usecase/init_refresh_client_secret.go"

type InitRefreshClientSecretDBO interface {
	SelectClientsBy(clientID string) (*entity.Client, *errorkit.DetailedError)
	InsertClientRegistrations(*entity.ClientRegistration) *errorkit.DetailedError
}

type InitRefreshClientSecret struct {
	detailedErrDescGen errorkit.ErrDescGenerator
	dbo                InitRefreshClientSecretDBO
}

func (irc InitRefreshClientSecret) Initiate(initClientID, clientID, authzToken string) (*entity.ClientRegistration, *errorkit.DetailedError) {
	client, detailedErr := irc.dbo.SelectClientsBy(clientID)
	if detailedErr != nil {
		return nil, detailedErr
	}

	ok, detailedErr := client.Authenticate(authzToken, irc.detailedErrDescGen)
	if !ok {
		return nil, detailedErr
	}

	expired, _ := client.IsExpired(irc.detailedErrDescGen)
	if !expired {
		return nil, nil
	}

	clientRegis, detailedErr := entity.NewClientRegistration(irc.detailedErrDescGen)
	if detailedErr != nil {
		return nil, detailedErr
	}
	clientRegis.InitClientID = initClientID

	detailedErr = irc.dbo.InsertClientRegistrations(clientRegis)
	if detailedErr != nil {
		return nil, detailedErr
	}

	return clientRegis, nil
}

func NewInitRefreshClientSecret(errDescGen errorkit.ErrDescGenerator, dbo InitRefreshClientSecretDBO) InitRefreshClientSecret {
	return InitRefreshClientSecret{errDescGen, dbo}
}
