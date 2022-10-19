package adapter

import "fmt"

const (
	OtherErrPasswordConstraint = iota
	OtherErrConfirmPasswordNoMatch
)

func GenerateOtherErrDesc(otherErrConst uint, args ...string) string {
	switch otherErrConst {
	case OtherErrPasswordConstraint:
		if len(args) == 1 {
			return fmt.Sprintf("%s must have a minimum 6 characters and maximum 32 character which consists of 1 lowercase, 1 uppercase, and 1 number", args[0])
		} else {
			return "a password must have a minimum 6 characters and maximum 32 character which consists of 1 lowercase, 1 uppercase, and 1 number"
		}
	case OtherErrConfirmPasswordNoMatch:
		return "confirm password doesn't match"
	default:
		return ""
	}
}
