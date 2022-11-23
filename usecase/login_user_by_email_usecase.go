package usecase

import (
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type LoginUserByEmailDBO interface {
	SelectUserWithRelBy(emailAddrXORUsername string, isEmail bool) (*entity.UserWithRel, *errorkit.DetailedError)
	SelectClientWithRelBy(clientID string) (*entity.ClientWithRel, *errorkit.DetailedError)
	InsertAuthzCodes(entity.AuthzCodeWithRel) (uint64, *errorkit.DetailedError)
	InsertScopes(entity.ScopeWithRel) *errorkit.DetailedError
}

type LoginUserByEmail struct {
	ErrDescGen       errorkit.ErrDescGenerator
	DBO              LoginUserByEmailDBO
	EmailXORUsername string
	ClientID         string
	RedirectURI      string
	IsEmail          bool
}

// must authenticate client secret and verify it's not expired
func (lue LoginUserByEmail) VerifyPassword(e2eEncPassword string) (bool, *errorkit.DetailedError) {
	// clientWithRel, detailedErr := lue.DBO.SelectClientWithRelBy(lue.ClientID)
	// if detailedErr != nil {
	// 	return false, detailedErr
	// }
	// user, detailedErr := lue.DBO.SelectUserWithRelBy(lue.EmailXORUsername, lue.IsEmail)
	// if detailedErr != nil {
	// 	return false, detailedErr
	// }
	return false, nil
}

func (lue LoginUserByEmail) GenerateAuthzCode() (string, *errorkit.DetailedError) {
	return "", nil
}

func NewLoginUserByEmail(errDescGen errorkit.ErrDescGenerator, dbo LoginUserByEmailDBO, emailXORUsername, clientID, redirectURI string, isEmail bool) LoginUserByEmail {
	return LoginUserByEmail{ErrDescGen: errDescGen, DBO: dbo, EmailXORUsername: emailXORUsername, ClientID: clientID, RedirectURI: redirectURI, IsEmail: isEmail}
}
