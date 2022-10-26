package entity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
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

func NewURLOneTimeToken(clientsID uint64, clientID string, storedClientSecretRaw []byte, host string, path string, query url.Values, fragment string, errDescGen errorkit.ErrDescGenerator) (*URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*URLOneTimeToken.GenerateURL", callTraceFileURLOneTimeToken)

	oneTimeToken, err := uuid.NewRandom()
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateRandomUUIDv4, errDescGen, "one time token")
	}
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrEd25519GenerateKeyPair, errDescGen, "url one time token")
	}
	signature := ed25519.Sign(sk, []byte(oneTimeToken.String()))

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
	query["one_time_token"] = []string{oneTimeToken.String()}
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
		OneTimeToken: oneTimeToken.String(),
		Signature:    base64.RawURLEncoding.EncodeToString(signature),
		ClientsID:    clientsID,
		Pk:           base64.RawURLEncoding.EncodeToString(pk),
		Sk:           base64.RawURLEncoding.EncodeToString(sk),
	}, nil
}
