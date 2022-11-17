package usecase

import (
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileInitRegisterClientUsecase = "/usecase/init_register_client_usecase.go"

type RegisterClientDBOperator interface {
	SelectCountClientRegistrationsBy(initClientID string) (int, *errorkit.DetailedError)
	InsertIgnore(clientRegistration *entity.ClientRegistration) *errorkit.DetailedError
	SelectCountClientsBy(initClientID string) (int, *errorkit.DetailedError)
}

func InitiateClientRegistration(initClientID string, dbo RegisterClientDBOperator, errDescGen errorkit.ErrDescGenerator) (*entity.ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#InitiateClientRegistration", callTraceFileInitRegisterClientUsecase)

	countedClient, detailedErr := dbo.SelectCountClientsBy(initClientID)
	if detailedErr != nil {
		return nil, detailedErr
	} else if countedClient > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrClientExists, errDescGen, "init_client_id")
	}

	countedRegis, detailedErr := dbo.SelectCountClientRegistrationsBy(initClientID)
	if detailedErr != nil {
		return nil, detailedErr
	} else if countedRegis == 0 {
		clientRegistration, detailedErr := entity.NewClientRegistration(errDescGen)
		if detailedErr != nil {
			return nil, detailedErr
		}
		clientRegistration.InitClientID = initClientID
		dbo.InsertIgnore(clientRegistration)
		return clientRegistration, nil
	} else {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrClientInitiatedRegister, errDescGen, initClientID)
	}
}
