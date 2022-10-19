package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileRegisterUserController = "/controller/register_user_controller.go"

type RegisterUser struct {
	dbo                  DBOperator
	htmlTemplateExecutor HTMLTemplateExecutor
}

func (ru RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.ServeHTTP", callTraceFileRegisterUserController)

	upv := restkit.URLParamValidation{RegexRules: map[string]uint{"client_id": regexkit.RegexUUIDV4}, Values: r.URL.Query()}

	regexErrMsgs, ok := upv.Validate(adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{FlattenRegexNoMatchMsgs: adapter.FlattenErrMessages(regexErrMsgs)}, w)
		return
	}

	registerUserRequest, detailedErr := adapter.ReadRequesBody[entity.RegisterUserRequest](r, "register user")
	if errorkit.IsNotNilThenLog(detailedErr) {
		w.WriteHeader(http.StatusInternalServerError)
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{Errs: []error{detailedErr}}, w)
		return
	}

	regexErrMsgs = ru.validateRequest(r, registerUserRequest)
	if len(*regexErrMsgs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{FlattenRegexNoMatchMsgs: adapter.FlattenErrMessages(regexErrMsgs)}, w)
		return
	}

	hpv := restkit.HeaderParamValidation{
		RegexRules: map[string]uint{"Csrf-Token": regexkit.RegexNotEmpty},
		Header:     r.Header,
	}

	regexErrMsgs, ok = hpv.Validate(adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{FlattenRegexNoMatchMsgs: adapter.FlattenErrMessages(regexErrMsgs)}, w)
		return
	}

	registerUserUsecase := usecase.RegisterUser{DBO: ru, ErrDescGen: errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc)}

	cookieKey := "csrf-token-hmac"
	csrfTokenHmac, err := r.Cookie(cookieKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{Errs: []error{errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrRetrieveCookie, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), cookieKey)}}, w)
		return
	}

	detailedErr = registerUserUsecase.VerifyCsrfTokenHmac(r.Header.Get("Csrf-Token"), csrfTokenHmac.Value)
	if errorkit.IsNotNilThenLog(detailedErr) {
		if detailedErr.Flow {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{Errs: []error{detailedErr}}, w)
		return
	}

	_, detailedErr = registerUserUsecase.CreateAndInsert(registerUserRequest.EmailAddress, registerUserRequest.Username, &registerUserRequest.Password, r.URL.Query().Get("client_id"))
	if errorkit.IsNotNilThenLog(detailedErr) {
		if detailedErr.Flow {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		ru.htmlTemplateExecutor.ExecuteTemplate("error_tmpl", "static/error.html", entity.ResponseBodyTemplate{Errs: []error{detailedErr}}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ru RegisterUser) validateRequest(r *http.Request, registerUserRequest *entity.RegisterUserRequest) *map[string][]string {
	var regexErrMsgs map[string][]string = make(map[string][]string)

	emailAddressOk := regexkit.RegexpCompiled[regexkit.RegexEmail].Match([]byte(registerUserRequest.EmailAddress))
	if !emailAddressOk {
		regexErrMsgs["email_address"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexEmail, "email_address")}
	}

	validatePassword("password", registerUserRequest.Password, &regexErrMsgs)
	validatePassword("confirm_password", registerUserRequest.ConfirmPassword, &regexErrMsgs)

	if registerUserRequest.Password != registerUserRequest.ConfirmPassword {
		regexErrMsgs["confirm_password"] = []string{adapter.GenerateOtherErrDesc(adapter.OtherErrConfirmPasswordNoMatch)}
	}

	return &regexErrMsgs
}

func validatePassword(key string, password string, regexErrMsgs *map[string][]string) {
	passwordOk := regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(password))
	if !passwordOk {
		(*regexErrMsgs)[key] = []string{adapter.RegexErrDescGen(regexkit.RegexNotEmpty)}
	} else {

		var numericFound bool
		var upperCaseFound bool
		var lowerCaseFound bool
		if len(password) < 6 {
			(*regexErrMsgs)[key] = []string{adapter.RegexErrDescGen(adapter.OtherErrPasswordConstraint, key)}
		} else {
			for idx := range password {
				if numericFound && upperCaseFound && lowerCaseFound {
					break
				}
				if password[idx] >= 48 && password[idx] <= 57 && !numericFound {
					numericFound = true
				} else if password[idx] >= 65 && password[idx] <= 90 && !upperCaseFound {
					upperCaseFound = true
				} else if password[idx] >= 97 && password[idx] <= 122 && !lowerCaseFound {
					lowerCaseFound = true
				}
			}
			if !numericFound || !upperCaseFound || !lowerCaseFound {
				(*regexErrMsgs)[key] = []string{adapter.RegexErrDescGen(adapter.OtherErrPasswordConstraint, key)}
			}
		}
	}
}

func (ru RegisterUser) SelectClientsIDBy(clientID string) (uint64, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.SelectClientsIDBy", callTraceFileRegisterUserController)

	row, err := ru.dbo.QueryRow("SELECT id FROM clients WHERE client_id = ?", clientID)
	if detailedErr := handleSelectDBNoRowsErr(err, callTraceFunc, "clients.id", "client_id"); detailedErr != nil {
		return 0, detailedErr
	}

	var clientsID uint64
	if scanErr := row.Scan(&clientsID); scanErr != nil {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc))
	}

	return clientsID, nil
}

func (ru RegisterUser) SelectCountUserBy(email, username string) (uint, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.SelectCountUserBy", callTraceFileRegisterUserController)

	row, err := ru.dbo.QueryRow("SELECT COUNT(u.id) FROM users u JOIN usernames un ON un.users_id = u.id WHERE u.email = ? AND un.username= ?", email, username)
	if detailedErr := handleSelectDBNoRowsErr(err, callTraceFunc, "COUNT(users.id)", "email", "username"); detailedErr != nil {
		return 0, detailedErr
	}

	var countedUsers uint
	if scanErr := row.Scan(&countedUsers); scanErr != nil {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "COUNT(users.id)", "countedUsers")
	}

	return countedUsers, nil
}

func (ru RegisterUser) SelectUsernameUnqNumBy(username string) (uint16, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.SelectUsernameUnqNumBy", callTraceFileRegisterUserController)

	row, err := ru.dbo.QueryRow("SELECT unq_num FROM usernames WHERE username = ?", username)
	if detailedErr := handleSelectDBNoRowsErr(err, callTraceFunc, "usernames.unq_num", "username"); detailedErr != nil {
		return 0, detailedErr
	}

	var unqNum uint16
	if scanErr := row.Scan(&unqNum); scanErr != nil {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, scanErr, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "usernames.unq_num", "unqNum")
	}

	return unqNum, nil
}

func (ru RegisterUser) InsertUser(user *entity.User) (uint64, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.InsertUser", callTraceFileRegisterUserController)

	commandResult, err := ru.dbo.Command(fmt.Sprintf("INSERT INTO users (email, password, clients_id) VALUES %s", sqlkit.GeneratePlaceHolder(3)), user.Email, user.Password, user.ClientsID)
	if err != nil {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "users")
	}

	lastInsertedID, err := commandResult.LastInsertId()
	if err != nil {
		return 0, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBLastInsertId, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "users")
	}

	return uint64(lastInsertedID), nil
}

func (ru RegisterUser) InsertUsername(username *entity.Username) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#RegisterUser.InsertUsername", callTraceFileRegisterUserController)

	_, err := ru.dbo.Command("INSERT INTO usernames (username, unq_num, users_id) VALUES (?, ?, ?)", username.Username, username.UnqNum, username.UsersID)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "usernames")
	}
	return nil
}

func NewRegisterUser(db *sql.DB, htmlTemplateExecutor HTMLTemplateExecutor) RegisterUser {
	return RegisterUser{dbo: &sqlkit.DBOperation{DB: db}, htmlTemplateExecutor: htmlTemplateExecutor}
}
