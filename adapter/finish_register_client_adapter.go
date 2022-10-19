package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileFinishRegisterClient = "/adapter/finish_register_client_adapter.go"

func MakeFinishClientRegistrationResponseBody(regexNoMatchMsgs *map[string][]string, message string, errors []error, result *entity.FinishClientRegistrationResult, shared *entity.FinishClientRegistrationShared) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeFinishClientRegistrationResponse", callTraceFileFinishRegisterClient)

	fcrr := entity.FinishClientRegistrationResponse{FinishClientRegistrationShared: shared, ResponseBodyTemplate: entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexNoMatchMsgs, Message: message, Errs: errors}, FinishClientRegistrationResult: result}
	if result != nil {
		fcrr.FinishClientRegistrationResult = result
	}
	response, err := json.Marshal(fcrr)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), "generate url one time token response body template"))
		return nil
	}
	return response
}
