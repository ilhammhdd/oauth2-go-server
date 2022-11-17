package entity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
)

const callTraceFileURLOneTimeToken = "/entity/url_one_time_token.go"

type URLOneTimeToken struct {
	ID            uint64     `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	Pk            string     `json:"pk,omitempty"`
	Sk            string     `json:"sk,omitempty"`
	OneTimeToken  string     `json:"one_time_token,omitempty"`
	Signature     string     `json:"signature,omitempty"`
	URL           string     `json:"path,omitempty"`
	ClientsID     uint64     `json:"clients_id,omitempty"`
}

func (uott *URLOneTimeToken) VerifySignature(signatureInReq string, errDescGen errorkit.ErrDescGenerator) (bool, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*URLOneTimeToken.VerifySignature", callTraceFileURLOneTimeToken)

	pk, err := base64.RawURLEncoding.DecodeString(uott.Pk)
	if err != nil {
		return false, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Encoding, errDescGen, "url one time token pk")
	}

	if len(pk) == 0 {
		return false, nil
	}

	signatureInReqRaw, err := base64.RawURLEncoding.DecodeString(signatureInReq)
	if err != nil {
		return false, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Encoding, errDescGen, "signature in req")
	}

	verified := ed25519.Verify(pk, []byte(uott.OneTimeToken), signatureInReqRaw)
	return verified, nil
}

func NewURLOneTimeToken(clientsID uint64, clientID string, host string, path string, query url.Values, fragment string, errDescGen errorkit.ErrDescGenerator) (*URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#NewURLOneTimeToken", callTraceFileURLOneTimeToken)

	oneTimeToken := GenerateRandID()
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrEd25519GenerateKeyPair, errDescGen, "url one time token")
	}
	signature := ed25519.Sign(sk, []byte(oneTimeToken))

	// trim the /init subpath to get the desired absolute path
	var urlBuilder strings.Builder
	if EnvVars[HTTPSEnvVar].Value == "false" {
		urlBuilder.WriteString("http://")
	} else if EnvVars[HTTPSEnvVar].Value == "true" {
		urlBuilder.WriteString("https://")
	}

	urlBuilder.WriteString(host)
	urlBuilder.WriteString(strings.Replace(path, "/init", "", 1))

	urlBuilder.WriteString("?")
	query["one_time_token"] = []string{oneTimeToken}
	query["signature"] = []string{base64.RawURLEncoding.EncodeToString(signature)}
	lastEle := len(query) - 1
	count := 0
	for urlQueryName, urlQueryValues := range query {
		urlBuilder.WriteString(urlQueryName)
		urlBuilder.WriteString("=")
		urlQueryValuesLen := len(urlQueryValues)
		for idx, urlQueryValue := range urlQueryValues {
			urlBuilder.WriteString(urlQueryValue)
			if idx < urlQueryValuesLen-1 {
				urlBuilder.WriteString(",")
			}
		}
		if count < lastEle {
			urlBuilder.WriteString("&")
		}
		count++
	}

	if fragment != "" {
		urlBuilder.WriteString("#")
		urlBuilder.WriteString(fragment)
	}

	return &URLOneTimeToken{
		URL:          urlBuilder.String(),
		OneTimeToken: oneTimeToken,
		Signature:    base64.RawURLEncoding.EncodeToString(signature),
		ClientsID:    clientsID,
		Pk:           base64.RawURLEncoding.EncodeToString(pk),
		Sk:           base64.RawURLEncoding.EncodeToString(sk),
	}, nil
}
