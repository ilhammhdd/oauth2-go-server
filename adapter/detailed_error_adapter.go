package adapter

import (
	"fmt"
	"strings"

	"ilhammhdd.com/oauth2-go-server/entity"
)

func GenerateDetailedErrDesc(errDescConst uint, args ...string) string {
	switch errDescConst {
	case entity.ErrGenerateCryptoRand:
		if len(args) == 1 {
			return fmt.Sprintf("error generating crypto random for %s", args[0])
		} else {
			return "error generating crypto random"
		}
	case entity.ErrParseToInt:
		if len(args) == 1 {
			return fmt.Sprintf("failed to parse %s to int", args[0])
		} else {
			return "parse to int failed"
		}
	case entity.ErrX25519Mul:
		if len(args) == 2 {
			return fmt.Sprintf("x25519 multiplication of %s and %s failed", args[0], args[1])
		} else {
			return "x25519 multiplication failed"
		}
	case entity.FlowErrDataNotMatched:
		if len(args) > 0 {
			return fmt.Sprintf("the following data not matched: [%s]", strings.Join(args, ", "))
		} else {
			return "data not matched"
		}
	case entity.FlowErrRegisterSessionExpired:
		return "registration session expired"
	case entity.FlowErrNotFoundBy:
		if len(args) >= 2 {
			return fmt.Sprintf("%s not found by: [%s]", args[0], strings.Join(args[1:], ", "))
		} else {
			return "not found by given params"
		}
	case entity.ErrUnsetEnvVar:
		if len(args) > 0 {
			return fmt.Sprintf("failed to unset env vars: [%s]", strings.Join(args, ", "))
		} else {
			return "failed to unset env vars"
		}
	case entity.ErrListenAndServe:
		return "error while listen and serve http server"
	case entity.FlowErrNotZeroValue:
		if len(args) > 0 {
			return fmt.Sprintf("args not a zero value: [%s]", strings.Join(args, ", "))
		} else {
			return "not a zero value"
		}
	case entity.ErrJsonMarshal:
		if len(args) == 1 {
			return fmt.Sprintf("failed to json marshal %s", args[0])
		} else {
			return "failed to marshal json"
		}
	case entity.ErrJsonUnmarshal:
		if len(args) == 1 {
			return fmt.Sprintf("failed to json unmarshal %s", args[0])
		} else {
			return "failed to unmarshal json"
		}
	case entity.FlowErrClientInitiatedRegister:
		return "client has initiated registration"
	case entity.ErrDBInsert:
		if len(args) == 1 {
			return fmt.Sprintf("error inserting to table %s", args[0])
		} else {
			return "error inserting to DB"
		}
	case entity.ErrDBSelect:
		if len(args) == 1 {
			return fmt.Sprintf("error selecting %s from DB", args[0])
		} else {
			return "error selecting from DB"
		}
	case entity.ErrDBLastInsertId:
		if len(args) == 1 {
			return fmt.Sprintf("error getting last inserted id of %s", args[0])
		} else {
			return "error getting last inserted id"
		}
	case entity.ErrDBScan:
		if len(args) == 2 {
			return fmt.Sprintf("error scanning %s from DB to %s", args[0], args[1])
		} else {
			return "error scanning from DB"
		}
	case entity.ErrDBDelete:
		if len(args) == 1 {
			return fmt.Sprintf("error deleting %s from DB", args[0])
		} else {
			return "error deleting from DB"
		}
	case entity.ErrBase64Encoding:
		if len(args) == 1 {
			return fmt.Sprintf("error base64 encoding %s", args[0])
		} else {
			return "error base64 encoding"
		}
	case entity.ErrBase64Decoding:
		if len(args) == 1 {
			return fmt.Sprintf("error base64 decoding %s", args[0])
		} else {
			return "error base64 decoding"
		}
	case entity.ErrSql:
		return "error sql operation"
	case entity.ErrGenerateRandomUUIDv4:
		if len(args) == 1 {
			return fmt.Sprintf("error generating random UUID v4 for %s", args[0])
		} else {
			return "error generating random UUID v4"
		}
	case entity.FlowErrClientExists:
		if len(args) == 1 {
			return fmt.Sprintf("client with the same %s already exists", args[0])
		} else {
			return "client already exists"
		}
	case entity.ErrTemplateHTML:
		if len(args) == 2 {
			return fmt.Sprintf("error %s for %s html template", args[0], args[1])
		} else {
			return "error html template"
		}
	case entity.ErrRetrieveCookie:
		if len(args) == 1 {
			return fmt.Sprintf("error while retrieving %s cookie", args[0])
		} else {
			return "error while retrieving cookie"
		}
	case entity.FlowErrBearerAccessTokenNotFound:
		return "Bearer access token not found"
	case entity.FlowErrUnauthorizedBearerAccessToken:
		return "unauthorized Bearer access token"
	case entity.ErrRequiredColumnIsNil:
		if len(args) == 1 {
			return fmt.Sprintf("required column %s is nil", args[0])
		} else {
			return "required column is nil"
		}
	case entity.FlowErrExistsBasedOn:
		if len(args) > 1 {
			return fmt.Sprintf("%s already exists based on [%s]", args[0], strings.Join(args[1:], ", "))
		} else {
			return "already exists"
		}
	case entity.ErrEd25519GenerateKeyPair:
		if len(args) == 1 {
			return fmt.Sprintf("error while generating ed25519 key pair for %s", args[0])
		} else {
			return "error while generating ed25519 key pair"
		}
	case entity.ErrReadRequestBody:
		if len(args) == 1 {
			return fmt.Sprintf("error while reading http request of %s", args[0])
		} else {
			return "error while reading http request"
		}
	case entity.FlowErrInvalidCsrfToken:
		return "invalid csrf token based on its signature"
	default:
		return "error description constant undefined"
	}
}
