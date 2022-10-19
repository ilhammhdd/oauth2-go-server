package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileHttpResponseAdapter = "/adapter/http_response_adapter.go"

func MakeResponseTmplErrResponse(regexErrMsgs *map[string][]string, message string, errs []error, errDescArgs ...string) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeResponseTmplErrResponse", callTraceFileHttpResponseAdapter)
	responseBodyTemplate := entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexErrMsgs, Message: message, Errs: errs}
	jsonBody, err := json.Marshal(responseBodyTemplate)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), errDescArgs...))
		return nil
	}

	return jsonBody
}
