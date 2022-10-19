package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileHttpRequestAdapter = "/adapter/http_request_adapter.go"

func ReadRequesBody[T interface{}](r *http.Request, detailedErrArgs ...string) (*T, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#ReadRequesBody", callTraceFileHttpRequestAdapter)

	bodyData := make([]byte, r.ContentLength)
	r.Body.Read(bodyData)
	// defer r.Body.Close()

	var result T
	err := json.Unmarshal(bodyData, &result)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonUnmarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), detailedErrArgs...)
	}

	return &result, nil
}
