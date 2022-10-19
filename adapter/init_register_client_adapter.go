package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileInitRegisterClientAdapter = "/adapter/init_register_client_adapter.go"

type InitiateRegisterClientResponse struct {
	*entity.ClientRegistration
	entity.ResponseBodyTemplate
}

// strip template columns from ClientRegistration before generating the response json
func MakeRegisterInitiateClientResponseBody(regexNoMatchMsgs *map[string][]string, errs []error, crIn *entity.ClientRegistration) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeRegisterInitiateClientResponseBody", callTraceFileInitRegisterClientAdapter)
	if crIn != nil {
		crIn.ID = nil
		crIn.CreatedAt = nil
		crIn.UpdatedAt = nil
		crIn.SoftDeletedAt = nil
		crIn.ServerSK = ""
	}

	response, err := json.Marshal(InitiateRegisterClientResponse{crIn, entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexNoMatchMsgs, Message: "", Errs: errs}})
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), "InitiateRegisterClientResponse"))
	}

	return response
}
