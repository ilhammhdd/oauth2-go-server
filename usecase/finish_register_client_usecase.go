package usecase

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileFinishRegisterClientUsecase = "/usecase/finish_register_client_usecase.go"

type FinishClientRegistrationDBOperator interface {
	SelectClientRegistrationsBy(initClientID string) (*entity.ClientRegistration, *errorkit.DetailedError)
	InsertClientWithRel(*entity.ClientWithRel) *errorkit.DetailedError
	DeleteClientRegistration(initClientID string) *errorkit.DetailedError
}

func FinishClientRegistration(fcrr *entity.FinishClientRegistrationRequest, authzToken string, scopeWithRels []*entity.ScopeWithRel, dbo FinishClientRegistrationDBOperator, errDescGen errorkit.ErrDescGenerator) (*entity.FinishClientRegistrationResult, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#FinishClientRegistration", callTraceFileFinishRegisterClientUsecase)

	defer dbo.DeleteClientRegistration(fcrr.InitClientID)

	cr, detailedErr := dbo.SelectClientRegistrationsBy(fcrr.InitClientID)
	if detailedErr != nil {
		return nil, detailedErr
	}

	fcrr.ServerSK = cr.ServerSK
	serverPkStored, err := base64.RawURLEncoding.DecodeString(cr.ServerPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "server_pk")
	}
	basepointStored, err := base64.RawURLEncoding.DecodeString(cr.Basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "basepoint")
	}

	serverPK, err := base64.RawURLEncoding.DecodeString(fcrr.ServerPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "stored server_pk")
	}
	basepoint, err := base64.RawURLEncoding.DecodeString(fcrr.Basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, errDescGen, "stored basepoint")
	}

	var errsDataNotMatched []string
	if !bytes.Equal(serverPK, serverPkStored) {
		errsDataNotMatched = append(errsDataNotMatched, "server_pk")
	}
	if !bytes.Equal(basepoint, basepointStored) {
		errsDataNotMatched = append(errsDataNotMatched, "basepoint")
	}
	if !fcrr.SessionExpiredAt.Equal(cr.SessionExpiredAt) {
		errsDataNotMatched = append(errsDataNotMatched, "session_expired_at")
	}
	if len(errsDataNotMatched) > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrDataNotMatched, errDescGen, errsDataNotMatched...)
	}

	sqlkit.TimeSubStripNano(fcrr.SessionExpiredAt, sqlkit.TimeNowUTCStripNano())
	now := sqlkit.TimeNowUTCStripNano()
	if now.Equal(fcrr.SessionExpiredAt) || now.After(fcrr.SessionExpiredAt) {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrRegisterSessionExpired, errDescGen)
	}

	clientID := entity.GenerateRandID()
	clientIDIssuedAt := sqlkit.TimeNowUTCStripNano()

	clientSecretRaw, detailedErr := fcrr.ClientRegistration.CalculateClientSecret(fcrr.ClientPK, errDescGen)
	if detailedErr != nil {
		return nil, detailedErr
	}
	clientSecretChecksumRaw := blake2b.Sum256(clientSecretRaw)
	if authzToken != base64.RawURLEncoding.EncodeToString(clientSecretChecksumRaw[:]) {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrUnauthorizedBearerAuthzToken, errDescGen)
	}

	clientSecretExpiredAt, detailedErr := entity.DetermineClientSecretExpiredAt(errDescGen)
	if detailedErr != nil {
		return nil, detailedErr
	}

	var redirectURIs []*entity.RedirectUri = make([]*entity.RedirectUri, len(fcrr.RedirectURIs))
	for idx := range fcrr.RedirectURIs {
		redirectURIs[idx] = &entity.RedirectUri{URI: fcrr.RedirectURIs[idx]}
	}
	var contacts []*entity.Contact = make([]*entity.Contact, len(fcrr.Contacts))
	for idx := range fcrr.Contacts {
		contacts[idx] = &entity.Contact{Contact: fcrr.Contacts[idx]}
	}

	if insertErr := dbo.InsertClientWithRel(&entity.ClientWithRel{
		Client: entity.Client{
			GrantTypes:              []byte(strings.Join(fcrr.GrantTypes, ",")),
			ResponseTypes:           []byte(strings.Join(fcrr.ResponseTypes, ",")),
			TokenEndpointAuthMethod: fcrr.TokenEndpointAuthMethod,
			ClientName:              fcrr.ClientName,
			ClientURI:               fcrr.ClientURI,
			LogoURI:                 fcrr.LogoURI,
			TosURI:                  fcrr.TosURI,
			PolicyURI:               fcrr.PolicyURI,
			SoftwareID:              fcrr.SoftwareID,
			SoftwareVersion:         fcrr.SoftwareVersion,
			InitClientID:            fcrr.InitClientID,
			ClientID:                clientID,
			ClientIDIssuedAt:        clientIDIssuedAt,
			ClientSecret:            base64.RawURLEncoding.EncodeToString(clientSecretRaw),
			ClientSecretExpiredAt:   clientSecretExpiredAt,
		},
		RedirectURIs:  redirectURIs,
		Contacts:      contacts,
		ScopesWithRel: scopeWithRels,
	}); insertErr != nil {
		return nil, insertErr
	}

	return &entity.FinishClientRegistrationResult{
		ClientID:              string(clientID),
		ClientIDIssuedAt:      clientIDIssuedAt,
		ClientSecretExpiredAt: clientSecretExpiredAt,
	}, nil
}
