package entity

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/curve25519"
)

const callTraceFileInitRegisterClientEntity = "/entity/init_register_client_entity.go"

type ClientRegistration struct {
	ID                   *uint64                   `json:"id,omitempty"`
	CreatedAt            *time.Time                `json:"created_at,omitempty"`
	UpdatedAt            *time.Time                `json:"updated_at,omitempty"`
	SoftDeletedAt        *time.Time                `json:"soft_deleted_at,omitempty"`
	InitClientIDChecksum string                    `json:"init_client_id_checksum,omitempty"`
	Basepoint            string                    `json:"basepoint,omitempty"`
	ServerSK             string                    `json:"server_sk,omitempty"`
	ServerPK             string                    `json:"server_pk,omitempty"`
	SessionExpiredAt     time.Time                 `json:"session_expired_at,omitempty"`
	errDescGen           errorkit.ErrDescGenerator `json:"-"`
}

func (crt *ClientRegistration) SetSessionExpiredAt() *errorkit.DetailedError {
	var callTraceFunc string = fmt.Sprintf("%s#ClientRegistration.SetSessionExpiredAt", callTraceFileInitRegisterClientEntity)
	// assuming the SessionExpiredAt has not been initialized with zero value
	if crt.SessionExpiredAt.Format(time.RFC3339Nano) != "0001-01-01T00:00:00Z" {
		return errorkit.NewDetailedError(true, callTraceFunc, nil, FlowErrNotZeroValue, crt.errDescGen, "session_expired_at")
	}

	var sessionExpiration int64
	sessionExpiration, err := strconv.ParseInt(EnvVars[ClientRegistrationSessionExpirationSecondsEnvVar].Value, 10, 64)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, ErrParseToInt, crt.errDescGen, "ClientRegistrationSessionExpirationSecondsEnvVar")
	} else {
		crt.SessionExpiredAt = time.Now().UTC().Add(time.Duration(sessionExpiration) * time.Second)
	}

	return nil
}

func NewInitiateClientRegistration(errDescGen errorkit.ErrDescGenerator) (*ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc string = fmt.Sprintf("%s#NewInitiateClientRegistration", callTraceFileInitRegisterClientEntity)
	basepoint, err := GenerateCryptoRand(curve25519.ScalarSize)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "basepoint")
	}

	secretKey, err := GenerateCryptoRand(curve25519.PointSize)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "secret_key")
	}

	publicKey, err := curve25519.X25519(secretKey, basepoint)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrX25519Mul, errDescGen, "basepoint", "secret_key")
	}

	return &ClientRegistration{Basepoint: base64.RawURLEncoding.EncodeToString(basepoint), ServerSK: base64.RawURLEncoding.EncodeToString(secretKey), ServerPK: base64.RawURLEncoding.EncodeToString(publicKey), errDescGen: errDescGen}, nil
}
