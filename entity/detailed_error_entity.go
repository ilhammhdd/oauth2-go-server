package entity

import "github.com/ilhammhdd/go-toolkit/errorkit"

const (
	NoErr uint = iota + errorkit.DetailedErrLastIota
	ErrGenerateCryptoRand
	ErrParseToInt
	ErrX25519Mul
	ErrUnsetEnvVar
	ErrListenAndServe
	ErrJsonMarshal
	ErrJsonUnmarshal
	ErrDBInsert
	ErrDBSelect
	ErrDBScan
	ErrDBDelete
	ErrDBLastInsertId
	ErrBase64Encoding
	ErrBase64Decoding
	ErrSql
	ErrGenerateRandomUUIDv4
	ErrTemplateHTML
	ErrRetrieveCookie
	ErrRequiredColumnIsNil
	ErrEd25519GenerateKeyPair
	ErrReadRequestBody
	ErrLastIota
)

const (
	FlowErrDataNotMatched uint = iota + ErrLastIota
	FlowErrRegisterSessionExpired
	FlowErrNotFoundBy
	FlowErrNotZeroValue
	FlowErrClientInitiatedRegister
	FlowErrClientExists
	FlowErrBearerAccessTokenNotFound
	FlowErrUnauthorizedBearerAccessToken
	FlowErrExistsBasedOn
	FlowErrInvalidCsrfToken
)
