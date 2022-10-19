package usecase

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileFinishRegisterClientUsecase = "/usecase/finish_register_client_usecase.go"

type FinishClientRegistrationDBOperator interface {
	SelectClientRegistrationBy(initClientIdChecksum string) (*entity.ClientRegistration, *errorkit.DetailedError)
	InsertClientWithRelations(*entity.ClientWithRelations) *errorkit.DetailedError
	DeleteInitClientRegistration(initClientIDChecksum string) *errorkit.DetailedError
}

func FinishClientRegistration(fcrr *entity.FinishClientRegistrationRequest, accessToken string, dbo FinishClientRegistrationDBOperator, errDescGen errorkit.ErrDescGenerator) (*entity.FinishClientRegistrationResult, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#FinishClientRegistration", callTraceFileFinishRegisterClientUsecase)

	initClientIDChecksum := fcrr.HashAndEncodeInitClientID()
	icr, detailedErr := dbo.SelectClientRegistrationBy(initClientIDChecksum)
	if errorkit.IsNotNilThenLog(detailedErr) {
		return nil, detailedErr
	}
	/* if err != nil && err == sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, err, entity.FlowErrNotFoundBy, errDescGen, "client_registrations", "init_client_id_checksum")
	} else if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrSql, errDescGen)
	} */

	fcrr.ServerSK = icr.ServerSK
	serverPkStored, err := base64.RawURLEncoding.DecodeString(icr.ServerPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "server_pk")
	}
	basepointStored, err := base64.RawURLEncoding.DecodeString(icr.Basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "basepoint")
	}

	icrServerPk, err := base64.RawURLEncoding.DecodeString(icr.ServerPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "stored server_pk")
	}
	icrBasepoint, err := base64.RawURLEncoding.DecodeString(icr.Basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "stored basepoint")
	}
	if !bytes.Equal(icrServerPk, serverPkStored) && bytes.Equal(icrBasepoint, basepointStored) && fcrr.SessionExpiredAt.Equal(icr.SessionExpiredAt) {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrDataNotMatched, errDescGen, "server_pk", "basepoint", "session_expired_at")
	}

	if fcrr.SessionExpiredAt.Sub(time.Now().UTC()).Nanoseconds() <= 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrRegisterSessionExpired, errDescGen)
	}

	clientID, err := uuid.NewRandom()
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrGenerateRandomUUIDv4, errDescGen, "client_id")
	}
	clientIDIssuedAt := time.Now().UTC()
	clientPk, err := base64.RawURLEncoding.DecodeString(fcrr.ClientPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "client_pk")
	}

	clientSecret, detailedErr := fcrr.CalculateClientSecret(clientPk, errDescGen)
	if notOk := errorkit.IsNotNilThenLog(detailedErr); notOk {
		return nil, detailedErr
	}
	clientSecretChecksum := blake2b.Sum256(clientSecret)
	if accessToken != base64.RawURLEncoding.EncodeToString(clientSecretChecksum[:]) {
		dbo.DeleteInitClientRegistration(initClientIDChecksum)
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrUnauthorizedBearerAccessToken, errDescGen)
	}

	clientSecretExpiredAt := time.Now().UTC().Add(30 * 24 * time.Hour)
	var redirectURIs []entity.RedirectUri = make([]entity.RedirectUri, len(fcrr.RedirectURIs))
	for idx := range fcrr.RedirectURIs {
		redirectURIs[idx] = entity.RedirectUri{URI: fcrr.RedirectURIs[idx]}
	}
	var contacts []entity.Contact = make([]entity.Contact, len(fcrr.Contacts))
	for idx := range fcrr.Contacts {
		contacts[idx] = entity.Contact{Contact: fcrr.Contacts[idx]}
	}

	if scanErr := dbo.InsertClientWithRelations(&entity.ClientWithRelations{
		Client: entity.Client{
			GrantTypes:              []byte(strings.Join(fcrr.GrantTypes, ",")),
			ResponseTypes:           []byte(strings.Join(fcrr.ResponseTypes, ",")),
			TokenEndpointAuthMethod: fcrr.TokenEndpointAuthMethod,
			ClientName:              fcrr.ClientName,
			ClientURI:               fcrr.ClientURI,
			LogoURI:                 fcrr.LogoURI,
			Scope:                   fcrr.Scope,
			TosURI:                  fcrr.TosURI,
			PolicyURI:               fcrr.PolicyURI,
			SoftwareID:              fcrr.SoftwareID,
			SoftwareVersion:         fcrr.SoftwareVersion,
			InitClientIDChecksum:    initClientIDChecksum,
			ClientID:                clientID.String(),
			ClientIDIssuedAt:        clientIDIssuedAt,
			ClientSecret:            base64.RawURLEncoding.EncodeToString(clientSecret),
			ClientSecretExpiredAt:   clientSecretExpiredAt,
		},
		RedirectURIs: redirectURIs,
		Contacts:     contacts,
	}); scanErr == nil {
		dbo.DeleteInitClientRegistration(initClientIDChecksum)
	}

	return &entity.FinishClientRegistrationResult{
		ClientID:              string(clientID.String()),
		ClientIDIssuedAt:      clientIDIssuedAt,
		ClientSecretExpiredAt: clientSecretExpiredAt,
	}, nil
}
