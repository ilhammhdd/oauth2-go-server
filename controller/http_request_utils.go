package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
)

func IsAccessTokenExists(w *http.ResponseWriter, r *http.Request, callTraceFunc string) (accessToken string, ok bool, statusCode int, responseBody []byte) {
	accessToken = ""
	ok = false
	statusCode = http.StatusOK
	responseBody = nil

	bearerAccessToken := r.Header.Get("Authorization")
	if strings.Contains(bearerAccessToken, "Bearer ") {
		if split := strings.Split(bearerAccessToken, " "); len(split) == 2 {
			accessToken = split[1]
			ok = true
		}
	}

	if accessToken == "" && !ok {
		response := entity.ResponseBodyTemplate{Errs: []error{errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrBearerAccessTokenNotFound, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc))}}
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
	return
}
