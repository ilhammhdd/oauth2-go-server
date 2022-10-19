package usecase

import (
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileInitRegisterClientUsecase = "/usecase/init_register_client_usecase.go"

type RegisterClientDBOperator interface {
	SelectCountBy(initClientIdChecksum string) (int, *errorkit.DetailedError)
	InsertIgnoreDBO(clientRegistration *entity.ClientRegistration) *errorkit.DetailedError
	SelectCountClientsBy(initClientIDChecksum string) (int, *errorkit.DetailedError)
}

func InitiateClientRegistration(initClientIDChecksum string, dbo RegisterClientDBOperator, errDescGen errorkit.ErrDescGenerator) (*entity.ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#InitiateClientRegistration", callTraceFileInitRegisterClientUsecase)

	countedClient, detailedErr := dbo.SelectCountClientsBy(initClientIDChecksum)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	} else if countedClient > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrClientExists, errDescGen, "init_client_id_checksum")
	}

	countedRegis, detailedErr := dbo.SelectCountBy(initClientIDChecksum)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	} else if countedRegis == 0 {
		clientRegistration, detailedErr := entity.NewInitiateClientRegistration(errDescGen)
		if errorkit.IsNotNilThenLog(detailedErr) {
			return nil, detailedErr
		}
		clientRegistration.InitClientIDChecksum = initClientIDChecksum
		dbo.InsertIgnoreDBO(clientRegistration)
		return clientRegistration, nil
	} else {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrClientInitiatedRegister, errDescGen, initClientIDChecksum)
	}
}
