package adapter

import (
	"fmt"

	"github.com/ilhammhdd/go-toolkit/regexkit"
)

const (
	RegexClientRegisterType uint = iota + regexkit.LastRegexIota
	RegexAppliationJson
	RegexTokenEndpointAuthMethod
	RegexGrantTypes
	RegexResponseTypes
	RegexBearerAuthzToken
	RegexBase64RawURL
	RegexRandomID
)

var Regex map[uint]string = map[uint]string{
	RegexClientRegisterType:      `^(init|refresh)$`,
	RegexAppliationJson:          `^(application/json)$`,
	RegexTokenEndpointAuthMethod: `^(none|client_secret_post|client_secret_basic|client_secret_bearer)$`,
	RegexGrantTypes:              `^(authorization_code|implicit|password|client_credentials|refresh_token)$`,
	RegexResponseTypes:           `^(code|token)$`,
	RegexBearerAuthzToken:        `^(Bearer [0-9a-zA-Z_-]{32,})$`,
	RegexBase64RawURL:            `^([0-9a-zA-Z_-]*)$`,
	RegexRandomID:                `^[a-zA-Z0-9_\-]{8}$`,
}

func FlattenMapSliceString(mapErrMessage *map[string][]string) []string {
	var flatten []string
	for key, val := range *mapErrMessage {
		for idx, errMessage := range val {
			if len(val) == 1 {
				flatten = append(flatten, fmt.Sprintf("%s: %s", key, errMessage))
			} else {
				flatten = append(flatten, fmt.Sprintf("%s[%d]: %s", key, idx, errMessage))
			}
		}
	}
	return flatten
}
