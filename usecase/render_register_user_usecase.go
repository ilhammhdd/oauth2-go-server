package usecase

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

// const callTraceFileRegisterUserUsecase = "/usecase/register_user_usecase.go"

type RenderRegisterUserDBOperator interface {
	DeleteURLOneTimeToken(oneTimeToken string, signature string) *errorkit.DetailedError
	SelectURLOneTimeToken(clientID string, oneTimeToken string, signature string) (*entity.URLOneTimeToken, *errorkit.DetailedError)
}

type RenderRegisterUser struct {
	ClientID     string
	OneTimeToken string
	ReqSignature string
	DBO          RenderRegisterUserDBOperator
	ErrDescGen   errorkit.ErrDescGenerator
}

func (rru *RenderRegisterUser) VerifySignature() (bool, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#RenderRegisterUser.VerifySignature", callTraceFileRegisterUserUsecase)

	signature, err := base64.RawURLEncoding.DecodeString(rru.ReqSignature)
	if err != nil {
		return false, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, rru.ErrDescGen, "url one time token signature")
	}

	urlOneTimeToken, detailedErr := rru.DBO.SelectURLOneTimeToken(rru.ClientID, rru.OneTimeToken, rru.ReqSignature)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return false, detailedErr
	}

	pk, err := base64.RawURLEncoding.DecodeString(urlOneTimeToken.Pk)
	if err != nil {
		return false, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Encoding, rru.ErrDescGen, "url one time token pk")
	}

	verified := ed25519.Verify(pk, []byte(urlOneTimeToken.OneTimeToken), signature)
	if !verified {
		deleteErr := rru.DBO.DeleteURLOneTimeToken(urlOneTimeToken.OneTimeToken, rru.ReqSignature)
		if deleteErr != nil {
			return false, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBDelete, rru.ErrDescGen, "url_one_time_tokens")
		}
	}
	return verified, nil
}

func (rru *RenderRegisterUser) GenerateCsrfTokenAndHmac() (csrfToken, csrfTokenHmac string, detailedErr *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*RenderRegisterUser.GenerateCsrfTokenAndHmac",
		callTraceFileRegisterUserUsecase)

	csrfTokenRaw, err := entity.GenerateCryptoRand(32)
	if err != nil {
		return "", "", errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrGenerateCryptoRand, rru.ErrDescGen, "csrf-token")
	}
	csrfToken = base64.RawURLEncoding.EncodeToString(csrfTokenRaw)

	csrfTokenHmacRaw := hmac.New(sha256.New, []byte(*entity.EphemeralCsrfTokenKey)).Sum(csrfTokenRaw)
	csrfTokenHmac = base64.RawURLEncoding.EncodeToString(csrfTokenHmacRaw)

	detailedErr = nil

	return
}
