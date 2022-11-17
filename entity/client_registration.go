package entity

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"golang.org/x/crypto/curve25519"
)

const callTraceFileClientRegistrationEntity = "/entity/client_registration.go"

type ClientRegistration struct {
	ID                 *uint64                   `json:"id,omitempty"`
	CreatedAt          *time.Time                `json:"created_at,omitempty"`
	UpdatedAt          *time.Time                `json:"updated_at,omitempty"`
	SoftDeletedAt      *time.Time                `json:"soft_deleted_at,omitempty"`
	InitClientID       string                    `json:"init_client_id,omitempty"`
	Basepoint          string                    `json:"basepoint,omitempty"`
	ServerSK           string                    `json:"server_sk,omitempty"`
	ServerPK           string                    `json:"server_pk,omitempty"`
	SessionExpiredAt   time.Time                 `json:"session_expired_at,omitempty"`
	detailedErrDescGen errorkit.ErrDescGenerator `json:"-"`
}

// Don't forget to call this func when inserting to persistent storage
func (cr *ClientRegistration) SetSessionExpiredAt() *errorkit.DetailedError {
	var callTraceFunc string = fmt.Sprintf("%s#ClientRegistration.SetSessionExpiredAt", callTraceFileClientRegistrationEntity)
	if !cr.SessionExpiredAt.IsZero() {
		return errorkit.NewDetailedError(true, callTraceFunc, nil, FlowErrNotZeroValue, cr.detailedErrDescGen, "session_expired_at")
	}

	var sessionExpiration int64
	sessionExpiration, err := strconv.ParseInt(EnvVars[ClientRegistrationSessionExpirationSecondsEnvVar].Value, 10, 64)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, ErrParseToInt, cr.detailedErrDescGen, "ClientRegistrationSessionExpirationSecondsEnvVar")
	} else {
		cr.SessionExpiredAt = sqlkit.TimeAddStripNano(sqlkit.TimeNowUTCStripNano(), time.Duration(sessionExpiration)*time.Second)
	}

	return nil
}

func (cr *ClientRegistration) CalculateClientSecret(clientPK string, errDescGen errorkit.ErrDescGenerator) ([]byte, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*ClientRegistration.CalculateSharedKey", callTraceFileClientRegistrationEntity)

	serverSKRaw, err := base64.RawURLEncoding.DecodeString(cr.ServerSK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Decoding, errDescGen, "client_registrations.server_sk")
	}

	clientPKRaw, err := base64.RawURLEncoding.DecodeString(clientPK)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrBase64Decoding, errDescGen, "client_pk")
	}

	var zeroValuesErrDescs []string
	if len(serverSKRaw) == 0 {
		zeroValuesErrDescs = append(zeroValuesErrDescs, "server_sk")
	}
	if len(clientPKRaw) == 0 {
		zeroValuesErrDescs = append(zeroValuesErrDescs, "client_pk")
	}
	if len(zeroValuesErrDescs) > 0 {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, nil, FlowErrZeroValue, cr.detailedErrDescGen, zeroValuesErrDescs...)
	}

	clientSecretRaw, err := curve25519.X25519(serverSKRaw, clientPKRaw)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrX25519Mul, errDescGen, "client_registrations.server_sk", "client_pk")
	}

	return clientSecretRaw, nil
}

func NewClientRegistration(errDescGen errorkit.ErrDescGenerator) (*ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc string = fmt.Sprintf("%s#NewInitiateClientRegistration", callTraceFileClientRegistrationEntity)
	basepoint, err := GenerateCryptoRand(curve25519.ScalarSize)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "basepoint")
	}

	serverSKRaw, err := GenerateCryptoRand(curve25519.PointSize)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "secret_key")
	}

	serverPKRaw, err := curve25519.X25519(serverSKRaw, basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrX25519Mul, errDescGen, "basepoint", "secret_key")
	}

	return &ClientRegistration{Basepoint: base64.RawURLEncoding.EncodeToString(basepoint), ServerSK: base64.RawURLEncoding.EncodeToString(serverSKRaw), ServerPK: base64.RawURLEncoding.EncodeToString(serverPKRaw), detailedErrDescGen: errDescGen}, nil
}
