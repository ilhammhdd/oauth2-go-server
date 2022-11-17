package entity

import (
	"time"
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
