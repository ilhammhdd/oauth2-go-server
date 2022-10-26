package entity

import (
	"encoding/json"
	"strings"
	"time"
)

type DBSet []byte

func (dbs DBSet) MarshalJSON() ([]byte, error) {
	dbsStr := string(dbs)
	dbsStrSplit := strings.Split(dbsStr, ",")
	return json.Marshal(dbsStrSplit)
}

func (dbs *DBSet) UnmarshalJSON(b []byte) error {
	var result DBSet
	for idx := range b {
		if b[idx] != '[' && b[idx] != ']' && b[idx] != '"' && b[idx] != '\n' && b[idx] != '\t' {
			result = append(result, b[idx])
		}
	}
	*dbs = result
	return nil
}

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
	Scope                   string     `json:"scope,omitempty"`
	TosURI                  string     `json:"tos_uri,omitempty"`
	PolicyURI               string     `json:"policy_uri,omitempty"`
	SoftwareID              string     `json:"software_id,omitempty"`
	SoftwareVersion         string     `json:"software_version,omitempty"`
	InitClientIDChecksum    string     `json:"init_client_id_checksum,omitempty"`
	ClientID                string     `json:"client_id,omitempty"`
	ClientIDIssuedAt        time.Time  `json:"client_id_issued_at,omitempty"`
	ClientSecret            string     `json:"client_secret,omitempty"`
	ClientSecretExpiredAt   time.Time  `json:"client_secret_expired_at,omitempty"`
}

type ClientWithRelations struct {
	Client
	RedirectURIs []RedirectUri `json:"redirect_uris,omitempty"`
	Contacts     []Contact     `json:"contacts,omitempty"`
}

type RedirectUri struct {
	ID            *uint64    `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	URI           string     `json:"uri,omitempty"`
	ClientsID     uint64     `json:"clients_id,omitempty"`
}

func FlattenRedirectURIsNonTemplateColumnsValue(redirectURIs *[]RedirectUri, clientsID uint64) []interface{} {
	var uris []interface{} = make([]interface{}, len(*redirectURIs)*2)
	for i := 0; i < len(*redirectURIs); i += 2 {
		uris[i*2] = (*redirectURIs)[i].URI
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
	ClientsId     uint64     `json:"clients_id,omitempty"`
}

func FlattenContactsNonTemplateColumnsValue(contacts *[]Contact, clientsID uint64) []interface{} {
	var uris []interface{} = make([]interface{}, len(*contacts)*2)
	for idx := range *contacts {
		uris[idx*2] = (*contacts)[idx].Contact
		uris[idx*2+1] = clientsID
	}
	return uris
}
