package adapter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileFinishRefreshClientSecret = "/adapter/finish_refresh_client_secret_adapter.go"

func MakeFinishRefreshClientSecretResponseBody(regexNoMatchMsgs *map[string][]string, message string, errors []error, clientID string, clientSecretExpiredAt time.Time) []byte {
	var callTraceFunc = fmt.Sprintf("%s#MakeFinishRefreshClientSecretResponseBody", callTraceFileFinishRefreshClientSecret)

	frcsr := entity.FinishRefreshClientSecretResponse{ClientSecretExpiredAt: clientSecretExpiredAt, ResponseBodyTemplate: entity.ResponseBodyTemplate{RegexNoMatchMsgs: regexNoMatchMsgs, Message: &message, Errs: errors}}

	response, err := json.Marshal(frcsr)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonMarshal, errorkit.ErrDescGeneratorFunc(GenerateDetailedErrDesc), "generate url one time token response body template"))
		return nil
	}
	return response
}
