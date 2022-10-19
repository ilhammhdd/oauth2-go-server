package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileGenerateURLOneTimeTokenAdapter = "/adapter/generate_url_one_time_token_adapter.go"

func MakeGenerateURLOneTimeTokenErrResponse(regexErrMsgs *map[string][]string, message string, errs []error) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeGenerateURLOneTimeTokenResponse", callTraceFileGenerateURLOneTimeTokenAdapter)
	responseBodyTemplate := entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexErrMsgs, Message: message, Errs: errs}
	jsonBody, err := json.Marshal(responseBodyTemplate)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), "generate url one time token response body template"))
		return nil
	}

	return jsonBody
}
