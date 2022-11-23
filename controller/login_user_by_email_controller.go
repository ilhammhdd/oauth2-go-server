package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileLoginUserByEmail = "/controller/login_user_by_email_controller.go"

type LoginUserByEmail struct {
	dbo DBOperator
}

func (lue LoginUserByEmail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upv := restkit.URLQueryValidation{RegexRules: map[string]uint{
		"client_id":     adapter.RegexRandomID,
		"redirect_uri":  regexkit.RegexURL,
		"scope":         regexkit.RegexNotEmpty,
		"response_type": adapter.RegexResponseTypeCode,
	}, Values: r.URL.Query()}
	regexNoMatchMsgs, ok := upv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNoMatchMsgs, "", nil))
		return
	}

	hpv := restkit.HeaderParamValidation{RegexRules: map[string]uint{
		"Authorization": adapter.RegexBearerAuthzToken,
	}, Header: r.Header}
	regexNoMatchMsgs, ok = hpv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNoMatchMsgs, "", nil))
		return
	}

	reqBody, detailedErr := adapter.ReadRequestBodyJSON[adapter.LoginUserByEmailRequest](r, "login user by email or username")
	if errorkit.IsNotNilThenLog(detailedErr) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isEmail, regexNoMatchMsgs := lue.validateRequestBody(*reqBody)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNoMatchMsgs, "", nil))
		return
	}

	uc := usecase.NewLoginUserByEmail(adapter.DetailedErrDescGen, lue, reqBody.EmailAddressUsername, r.URL.Query().Get("client_id"), r.URL.Query().Get("redirect_uri"), isEmail)
	uc.VerifyPassword(reqBody.Password)

}

func (lue LoginUserByEmail) validateRequestBody(reqBody adapter.LoginUserByEmailRequest) (isEmail bool, regexNoMatchMsgs *map[string][]string) {
	var noMatchMsgs map[string][]string = make(map[string][]string)

	if ok := regexkit.RegexpCompiled[regexkit.RegexEmail].Match([]byte(reqBody.EmailAddressUsername)); !ok {
		if ok := regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(reqBody.EmailAddressUsername)); !ok {
			noMatchMsgs["email_address_username"] = []string{"not a valid email or username"}
		}
	} else {
		isEmail = true
	}

	if ok := regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(reqBody.Password)); !ok {
		noMatchMsgs["password"] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	}

	return isEmail, &noMatchMsgs
}

func (lue LoginUserByEmail) SelectUserWithRelBy(emailAddrXORUsername string, isEmail bool) (*entity.UserWithRel, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#LoginUserByEmail.SelectUserWithRelBy", callTraceFileLoginUserByEmail)

	var stmt string
	if isEmail {
		stmt = "SELECT id,email,password FROM users WHERE email=?"
	} else {
		stmt = "SELECT u.id,u.email,u.password FROM users u JOIN usernames un ON un.users_id=u.id WHERE un.username=?"
	}

	var user entity.User
	userRow, err := lue.dbo.QueryRow(stmt, emailAddrXORUsername)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "users"); detailedErr != nil {
		return nil, detailedErr
	}

	if err := userRow.Scan(&user.ID, &user.Email, &user.Password); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen, "users", "user")
	}

	uppRow, err := lue.dbo.QueryRow("SELECT upp.rand_salt,upp.time,upp.memory,upp.threads,upp.keyLen FROM user_password_params upp JOIN users u ON upp.users_id=u.id WHERE u.id=?", user.ID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "user_password_params"); detailedErr != nil {
		return nil, detailedErr
	}

	var upp entity.UserPasswordParams
	if err := uppRow.Scan(&upp.RandSalt, &upp.Time, &upp.Memory, &upp.Threads, &upp.KeyLen); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen, "user_password_params", "upp")
	}

	return &entity.UserWithRel{User: &user, UserPasswordParams: &upp}, nil
}

func (lue LoginUserByEmail) SelectClientWithRelBy(clientID string) (*entity.ClientWithRel, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#LoginUserByEmail.SelectClientWithRelBy", callTraceFileLoginUserByEmail)

	row, err := lue.dbo.QueryRow("SELECT id,response_types,client_secret,client_secret_expired_at FROM clients WHERE client_id=?", clientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc); detailedErr != nil {
		return nil, detailedErr
	}

	var client entity.Client
	if err := row.Scan(&client.ID, &client.ResponseTypes, &client.ClientSecret, &client.ClientSecret, &client.ClientSecretExpiredAt); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen, "clients", "client")
	}

	rows, err := lue.dbo.Query("SELECT uri FROM redirect_uris WHERE clients_id=?", client.ID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc); detailedErr != nil {
		return nil, detailedErr
	}

	var redirectURIs []*entity.RedirectUri
	for rows.Next() {
		var redirectURI entity.RedirectUri
		if err := rows.Scan(&redirectURI.ID); err != nil && err != sql.ErrNoRows {
			return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, adapter.DetailedErrDescGen)
		}
		redirectURIs = append(redirectURIs, &redirectURI)
	}

	return &entity.ClientWithRel{Client: client, RedirectURIs: redirectURIs}, nil
}

func (lue LoginUserByEmail) InsertAuthzCodes(entity.AuthzCodeWithRel) (uint64, *errorkit.DetailedError) {
	return 0, nil
}

func (lue LoginUserByEmail) InsertScopes(entity.ScopeWithRel) *errorkit.DetailedError {
	return nil
}
