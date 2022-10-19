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
	RegexBearerAccessToken
	RegexBase64RawURL
)

var Regex map[uint]string = map[uint]string{
	RegexClientRegisterType:      `^(init|refresh)$`,
	RegexAppliationJson:          `^(application/json)$`,
	RegexTokenEndpointAuthMethod: `^(none|client_secret_post|client_secret_basic|client_secret_bearer)$`,
	RegexGrantTypes:              `^(authorization_code|implicit|password|client_credentials|refresh_token)$`,
	RegexResponseTypes:           `^(code|token)$`,
	RegexBearerAccessToken:       `^(Bearer [0-9a-zA-Z_-]{32,})$`,
	RegexBase64RawURL:            `^([0-9a-zA-Z_-]*)$`,
}

func FlattenErrMessages(mapErrMessage *map[string][]string) []string {
	var flatten []string
	for key, val := range *mapErrMessage {
		for idx, errMessage := range val {
			flatten = append(flatten, fmt.Sprintf("%s[%d] %s", key, idx, errMessage))
		}
	}
	return flatten
}
