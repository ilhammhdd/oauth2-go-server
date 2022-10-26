package adapter

import (
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
)

var RegexErrDescGen = errorkit.ErrDescGeneratorFunc(GenerateRegexErrDesc)

func generateDescOfArgs(desc string, args ...string) string {
	if len(args) == 2 {
		return fmt.Sprintf("%s[%s] %s", args[0], args[1], desc)
	} else if len(args) == 1 {
		return fmt.Sprintf("%s %s", args[0], desc)
	} else {
		return desc
	}
}

func GenerateRegexErrDesc(regexConst uint, args ...string) string {
	switch regexConst {
	case regexkit.RegexEmail:
		return generateDescOfArgs("is not a valid email", args...)
	case regexkit.RegexAlphanumeric:
		return generateDescOfArgs("is not a valid alphanumeric", args...)
	case regexkit.RegexNotEmpty:
		return generateDescOfArgs("is not allowed empty", args...)
	case regexkit.RegexURL:
		return generateDescOfArgs("is not a valid URL", args...)
	case regexkit.RegexJWT:
		return generateDescOfArgs("is not a valid JWT", args...)
	case regexkit.RegexNumber:
		return generateDescOfArgs("is not a valid number", args...)
	case regexkit.RegexLatitude:
		return generateDescOfArgs("is not a valid latitude", args...)
	case regexkit.RegexLongitude:
		return generateDescOfArgs("is not a valid longitude", args...)
	case regexkit.RegexUUIDV4:
		return generateDescOfArgs("is not a valid UUIDv4", args...)
	case regexkit.RegexCommonUnitOfLength:
		return generateDescOfArgs("is not a valid common unit of length", args...)
	case regexkit.RegexIPv4:
		return generateDescOfArgs("is not a valid IPv4", args...)
	case regexkit.RegexIPv4TCPPortRange:
		return generateDescOfArgs("is not a valid IPv4 TCP port range", args...)
	case RegexClientRegisterType:
		return generateDescOfArgs("is not a valid client register type", args...)
	case RegexAppliationJson:
		return generateDescOfArgs("is not an application/json", args...)
	case RegexTokenEndpointAuthMethod:
		return generateDescOfArgs("is not a valid token endpoint auth method", args...)
	case RegexGrantTypes:
		return generateDescOfArgs("is not a valid grant types", args...)
	case RegexResponseTypes:
		return generateDescOfArgs("is not a valid response types", args...)
	case regexkit.RegexDateTimeRFC3339:
		return generateDescOfArgs("is not a valid RFC3339 date time", args...)
	case RegexBearerAccessToken:
		return generateDescOfArgs("is not a valid Bearer access token", args...)
	case RegexBase64RawURL:
		return generateDescOfArgs("is not a valid base64 raw URL encoding", args...)
	default:
		return fmt.Sprintf("no error message for regex const [%d] found", regexConst)
	}
}
