package usecase

import (
	"fmt"
	"log"
	"net/url"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileGenerateURLOneTimeToken = "/usecase/generate_url_one_time_token.go"

// type SelectClientByClientIDResult struct {
// 	ID           uint64
// 	ClientID     string
// 	ClientSecret string
// }

type GenerateURLOneTimeTokenDBOperator interface {
	SelectClientsBy(clientID string) (*entity.Client, *errorkit.DetailedError)
	SelectCountURLOneTimeToken(clientsID uint64, url string) (uint32, *errorkit.DetailedError)
	InsertURLOneTimeToken(*entity.URLOneTimeToken) *errorkit.DetailedError
}

type GenerateURLOneTimeToken struct {
	ClientID   string
	AuthzToken string
	DBO        GenerateURLOneTimeTokenDBOperator
	ErrDescGen errorkit.ErrDescGenerator
}

func (guott *GenerateURLOneTimeToken) Generate(host string, path string, query url.Values, fragment string) (*entity.URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*GenerateURLOneTimeToken.Generate", callTraceFileGenerateURLOneTimeToken)

	client, detailedErr := guott.DBO.SelectClientsBy(guott.ClientID)
	if detailedErr != nil {
		return nil, detailedErr
	}

	ok, detailedErr := client.Authenticate(guott.AuthzToken, guott.ErrDescGen)
	if !ok {
		return nil, detailedErr
	}

	expired, detailedErr := client.IsExpired(guott.ErrDescGen)
	if expired {
		log.Println("generate url ott authz token expired")
		return nil, detailedErr
	}

	entityURLOneTimeToken, detailedErr := entity.NewURLOneTimeToken(*client.ID, client.ClientID, host, path, query, fragment, guott.ErrDescGen)
	if detailedErr != nil {
		return nil, detailedErr
	}

	existingURLOneTimeToken, detailedErr := guott.DBO.SelectCountURLOneTimeToken(*client.ID, entityURLOneTimeToken.URL)
	if detailedErr != nil {
		return nil, detailedErr
	}

	if existingURLOneTimeToken > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrExistsBasedOn, guott.ErrDescGen, "url_one_time_tokens", "clients.id", "url_one_time_tokens.url")
	}

	detailedErr = guott.DBO.InsertURLOneTimeToken(entityURLOneTimeToken)
	if detailedErr != nil {
		return nil, detailedErr
	}

	return entityURLOneTimeToken, nil
}
