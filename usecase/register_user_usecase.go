package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"ilhammhdd.com/oauth2-go-server/entity"
)

const callTraceFileRegisterUserUsecase = "/usecase/register_user_usecase.go"

type RegisterUserDBOperator interface {
	SelectClientsIDBy(clientID string) (uint64, *errorkit.DetailedError)
	SelectCountUserBy(email, username string) (uint, *errorkit.DetailedError)
	SelectUsernameUnqNumBy(username string) (uint16, *errorkit.DetailedError)
	InsertUser(user *entity.User) (uint64, *errorkit.DetailedError)
	InsertUsername(username *entity.Username) *errorkit.DetailedError
	InsertUserPasswordParams(userPasswordParams *entity.UserPasswordParams) *errorkit.DetailedError
}

type RegisterUser struct {
	OneTimeToken string
	Signature    string
	DBO          RegisterUserDBOperator
	ErrDescGen   errorkit.ErrDescGenerator
}

func (ru RegisterUser) VerifyCsrfTokenHmac(csrfToken, csrfTokenHmac string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.VerifyCsrfTokenHmac", callTraceFileRegisterUserUsecase)

	csrfTokenRaw, err := base64.RawURLEncoding.DecodeString(csrfToken)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, ru.ErrDescGen, "csrf-token")
	}

	calculatedCsrfTokenHmac := hmac.New(sha256.New, []byte(*entity.EphemeralCsrfTokenKey)).Sum(csrfTokenRaw)

	csrfTokenHmacRaw, err := base64.RawURLEncoding.DecodeString(csrfTokenHmac)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrBase64Decoding, ru.ErrDescGen, "csrf-token-hmac-regsiter")
	}

	if !hmac.Equal(calculatedCsrfTokenHmac, csrfTokenHmacRaw) {
		return errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrInvalidCsrfToken, ru.ErrDescGen)
	}

	return nil
}

func (ru RegisterUser) CreateAndInsert(email string, username string, plainPassword *string, clientID string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.CreateAndInsert", callTraceFileRegisterUserUsecase)

	countedUser, detailedErr := ru.DBO.SelectCountUserBy(email, username)
	if detailedErr != nil {
		return detailedErr
	}
	if countedUser > 0 {
		return errorkit.NewDetailedError(true, callTraceFunc, nil, entity.FlowErrExistsBasedOn, ru.ErrDescGen, "email", "username")
	}

	clientsID, detailedErr := ru.DBO.SelectClientsIDBy(clientID)
	if detailedErr != nil {
		return detailedErr
	}

	user, userPasswordParams, detailedErr := entity.NewUserWithPasswordParams(email, plainPassword, clientsID, ru.ErrDescGen)
	if detailedErr != nil {
		return detailedErr
	}

	usersID, detailedErr := ru.DBO.InsertUser(user)
	if detailedErr != nil {
		return detailedErr
	}

	userPasswordParams.UsersID = usersID
	detailedErr = ru.DBO.InsertUserPasswordParams(userPasswordParams)
	if detailedErr != nil {
		return detailedErr
	}

	usernameUnqNum, detailedErr := ru.DBO.SelectUsernameUnqNumBy(username)
	if detailedErr != nil && detailedErr != sql.ErrNoRows {
		return detailedErr
	}

	usernameEntity := entity.Username{Username: username, UsersID: usersID}
	if detailedErr == sql.ErrNoRows {
		usernameEntity.UnqNum = 1
	} else if usernameUnqNum > 0 {
		usernameEntity.UnqNum = usernameUnqNum + 1
	}

	detailedErr = ru.DBO.InsertUsername(&usernameEntity)
	if detailedErr != nil {
		return detailedErr
	}

	return nil
}
