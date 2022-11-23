package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileHttpResponseAdapter = "/adapter/http_response_adapter.go"

func MakeResponseTmplErrResponse(regexErrMsgs *map[string][]string, message string, detailedErrs ...*errorkit.DetailedError) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeResponseTmplErrResponse", callTraceFileHttpResponseAdapter)
	var responseBodyTemplate entity.ResponseBodyTemplate
	if regexErrMsgs != nil {
		responseBodyTemplate = entity.ResponseBodyTemplate{FlatRegexNoMatchMsgs: FlattenMapSliceString(regexErrMsgs), RegexNoMatchMsgs: regexErrMsgs, Message: &message, DetailedErrs: detailedErrs}
	} else {
		responseBodyTemplate = entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexErrMsgs, Message: &message, DetailedErrs: detailedErrs}
	}
	jsonBody, err := json.Marshal(responseBodyTemplate)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc)))
		return nil
	}

	return jsonBody
}
