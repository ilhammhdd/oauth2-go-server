package usecase

import (
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type UrlOneTimeTokenUsecaseDBO interface {
	SelectURLOneTimeToken(clientID string, oneTimeToken string, signature string) (*entity.URLOneTimeToken, *errorkit.DetailedError)
	DeleteURLOneTimeToken(oneTimeToken string, signature string) *errorkit.DetailedError
}

type UrlOneTimeTokenUsecase struct {
	ErrDescGen errorkit.ErrDescGenerator
	DBO        UrlOneTimeTokenUsecaseDBO
}

func (uottu *UrlOneTimeTokenUsecase) VerifySignature(clientID, oneTimeToken, signature string) (bool, *errorkit.DetailedError) {
	urlOneTimeToken, detailedErr := uottu.DBO.SelectURLOneTimeToken(clientID, oneTimeToken, signature)
	if detailedErr != nil {
		return false, detailedErr
	}

	verified, detailedErr := urlOneTimeToken.VerifySignature(signature, uottu.ErrDescGen)
	if detailedErr != nil {
		return false, detailedErr
	}

	detailedErr = uottu.DBO.DeleteURLOneTimeToken(oneTimeToken, signature)
	if detailedErr != nil {
		return verified, detailedErr
	}

	return verified, nil
}
