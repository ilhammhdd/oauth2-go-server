package entity

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"golang.org/x/crypto/blake2b"
)

const callTraceFileClient = "/entity/client.go"

type Client struct {
	ID                      *uint64    `json:"id,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt           *time.Time `json:"soft_deleted_at,omitempty"`
	GrantTypes              DBSet      `json:"grant_types,omitempty"`
	ResponseTypes           DBSet      `json:"response_types,omitempty"`
	TokenEndpointAuthMethod string     `json:"token_endpoint_auth_method,omitempty"`
	ClientName              string     `json:"client_name,omitempty"`
	ClientURI               string     `json:"client_uri,omitempty"`
	LogoURI                 string     `json:"logo_uri,omitempty"`
	TosURI                  string     `json:"tos_uri,omitempty"`
	PolicyURI               string     `json:"policy_uri,omitempty"`
	SoftwareID              string     `json:"software_id,omitempty"`
	SoftwareVersion         string     `json:"software_version,omitempty"`
	InitClientID            string     `json:"init_client_id,omitempty"`
	ClientID                string     `json:"client_id,omitempty"`
	ClientIDIssuedAt        time.Time  `json:"client_id_issued_at,omitempty"`
	ClientSecret            string     `json:"client_secret,omitempty"`
	ClientSecretExpiredAt   time.Time  `json:"client_secret_expired_at,omitempty"`
}

func (client *Client) Authenticate(inAuthzToken string, errDescGen errorkit.ErrDescGenerator) (bool, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*Client.Authenticate", callTraceFileClient)

	clientSecretRaw, err := base64.RawURLEncoding.DecodeString(client.ClientSecret)
	if err != nil {
		return false, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Decoding, errDescGen, "clients")
	}

	clientSecretChecksum := blake2b.Sum256(clientSecretRaw)
	storedClientAuthzToken := base64.RawURLEncoding.EncodeToString(clientSecretChecksum[:])
	if storedClientAuthzToken != inAuthzToken {
		return false, errorkit.NewDetailedError(true, callTraceFunc, nil, FlowErrUnauthorizedBearerAuthzToken, errDescGen)
	}

	return true, nil
}

func (client *Client) IsExpired(errDescGen errorkit.ErrDescGenerator) (bool, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*Client.IsExpired", callTraceFileClient)

	if sqlkit.TimeNowUTCStripNano().Before(client.ClientSecretExpiredAt) {
		return false, nil
	}

	return true, errorkit.NewDetailedError(true, callTraceFunc, nil, FlowErrBearerAuthzTokenExpired, errDescGen)
}

func DetermineClientSecretExpiredAt(errDescGen errorkit.ErrDescGenerator) (time.Time, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#DetermineClientSecretExpiredAt", callTraceFileClient)

	clientSecretExpirationSeconds, err := strconv.ParseInt(EnvVars[ClientSecretExpirationSecondsEnvVar].Value, 10, 64)
	if err != nil {
		return time.Time{}, errorkit.NewDetailedError(false, callTraceFunc, err, ErrParseToInt, errDescGen, ClientSecretExpirationSecondsEnvVar)
	}
	return sqlkit.TimeAddStripNano(sqlkit.TimeNowUTCStripNano(), time.Duration(clientSecretExpirationSeconds)*time.Second), nil
}

type ClientWithRel struct {
	Client
	RedirectURIs           []*RedirectUri           `json:"redirect_uris,omitempty"`
	Contacts               []*Contact               `json:"contacts,omitempty"`
	AccessXorRefreshTokens []*AccessXorRefreshToken `json:"access_xor_refresh_tokens"`
	ScopesWithRel          []*ScopeWithRel          `json:"scopes"`
}

type RedirectUri struct {
	ID            *uint64    `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	URI           string     `json:"uri,omitempty"`
	ClientsID     uint64     `json:"clients_id,omitempty"`
}

func FlattenRedirectURIsNonTemplateColumnsValue(redirectURIs []*RedirectUri, clientsID uint64) []interface{} {
	var uris []interface{} = make([]interface{}, len(redirectURIs)*2)
	for i := 0; i < len(redirectURIs); i += 2 {
		uris[i*2] = redirectURIs[i].URI
		uris[i*2+1] = clientsID
	}
	return uris
}

type Contact struct {
	Id            *uint64    `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	Contact       string     `json:"contact,omitempty"`
	ClientsID     uint64     `json:"clients_id,omitempty"`
}

func FlattenContactsNonTemplateColumnsValue(contacts []*Contact, clientsID uint64) []interface{} {
	var uris []interface{} = make([]interface{}, len(contacts)*2)
	for idx := range contacts {
		uris[idx*2] = contacts[idx].Contact
		uris[idx*2+1] = clientsID
	}
	return uris
}
