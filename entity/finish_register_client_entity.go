package entity

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/curve25519"
)

const callTraceFileFinisRegisterClientEntity = "/entity/finish_register_client_entity.go"

type FinishClientRegistrationShared struct {
	RedirectURIs            []string `json:"redirect_uris"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	ClientName              string   `json:"client_name"`
	ClientURI               string   `json:"client_uri"`
	LogoURI                 string   `json:"logo_uri"`
	Scope                   string   `json:"scope"`
	Contacts                []string `json:"contacts"`
	TosURI                  string   `json:"tos_uri"`
	PolicyURI               string   `json:"policy_uri"`
	SoftwareID              string   `json:"software_id"`
	SoftwareVersion         string   `json:"software_version"`
}

type FinishClientRegistrationRequest struct {
	InitClientID string `json:"init_client_id"`
	ClientPK     string `json:"client_pk"`
	FinishClientRegistrationShared
	ClientRegistration
}

type FinishClientRegistrationResult struct {
	ClientID              string    `json:"client_id"`
	ClientIDIssuedAt      time.Time `json:"client_id_issued_at"`
	ClientSecretExpiredAt time.Time `json:"client_secret_expired_at"`
}

type FinishClientRegistrationResponse struct {
	*FinishClientRegistrationResult
	ResponseBodyTemplate
	*FinishClientRegistrationShared
}

func (fcr *FinishClientRegistrationRequest) HashAndEncodeInitClientID() string {
	hash := blake2b.Sum256([]byte(fcr.InitClientID))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (fcr *FinishClientRegistrationRequest) CalculateClientSecret(clientPk []byte, errDescGen errorkit.ErrDescGenerator) ([]byte, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#FinishClientRegistrationRequest.CalculateClientSecret", callTraceFileFinisRegisterClientEntity)
	serverSk, err := base64.RawURLEncoding.DecodeString(fcr.ServerSK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Decoding, errDescGen, "server_sk")
	}

	clientSecret, err := curve25519.X25519(serverSk, clientPk)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrX25519Mul, errDescGen, "server_sk", "client_pk")
	}

	return clientSecret, nil
}

func (fcr *FinishClientRegistrationRequest) GenerateClientSecretExpiredAt() time.Time {
	time.Now().UTC().Add(30 * 24 * time.Hour)
	return time.Now()
}
