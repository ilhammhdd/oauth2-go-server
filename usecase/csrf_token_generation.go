package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileCsrfTokenGeneration = "/usecase/csrf_token_generation.go"

func GenerateCsrfTokenAndHmac(errDescGen errorkit.ErrDescGenerator) (csrfToken, csrfTokenHmac string, detailedErr *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#GenerateCsrfTokenAndHmac",
		callTraceFileCsrfTokenGeneration)

	csrfTokenRaw, err := entity.GenerateCryptoRand(32)
	if err != nil {
		return "", "", errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrGenerateCryptoRand, errDescGen, "csrf-token")
	}
	csrfToken = base64.RawURLEncoding.EncodeToString(csrfTokenRaw)

	csrfTokenHmacRaw := hmac.New(sha256.New, []byte(*entity.EphemeralCsrfTokenKey)).Sum(csrfTokenRaw)
	csrfTokenHmac = base64.RawURLEncoding.EncodeToString(csrfTokenHmacRaw)

	detailedErr = nil

	return
}
