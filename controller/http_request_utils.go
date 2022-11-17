package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

func IsAuthzTokenExists(prefix, authorizationHeader, callTraceFunc string) (string, int, []byte) {
	authzToken := ""
	statusCode := http.StatusOK
	var responseBody []byte = nil

	if strings.Contains(authorizationHeader, fmt.Sprintf("%s ", prefix)) {
		if split := strings.Split(authorizationHeader, " "); len(split) == 2 {
			authzToken = split[1]
		}
	}

	if authzToken == "" {
		response := entity.ResponseBodyTemplate{Errs: []error{errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrBearerAuthzTokenNotFound, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc))}}
		jsonBody, err := json.Marshal(response)
		if err != nil {
			detailedError := errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "response body template")
			errorkit.IsNotNilThenLog(detailedError)
			statusCode = http.StatusInternalServerError
		} else {
			statusCode = http.StatusUnauthorized
			responseBody = jsonBody
		}
	}

	return authzToken, statusCode, responseBody
}
