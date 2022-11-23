package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileFinishRegisterClientController = "/controller/finish_register_client_controller.go"

type FinishClientRegistration struct {
	dbo DBOperator
}

func (fcr *FinishClientRegistration) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var callTraceFunc = fmt.Sprintf("%s#*FinishClientRegistration.ServeHTTP", callTraceFileFinishRegisterClientController)
	w.Header().Set("Content-Type", "application/json")

	authzToken := fcr.getAuthzToken(r.Header.Get("Authorization"))

	fcrr, regexErrMsgs := fcr.validateAndParseRequest(r)
	if len(regexErrMsgs) > 0 {
		if _, ok := regexErrMsgs["Authorization"]; ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(nil)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeFinishClientRegistrationResponseBody(&regexErrMsgs, "", nil, nil, nil))
		return
	}

	scopeWithRels, detailedErrs := adapter.ParseScopeWithRel(fcrr.Scope, callTraceFunc)
	if len(detailedErrs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", detailedErrs...))
		return
	}
	result, detailedErr := usecase.FinishClientRegistration(fcrr, authzToken, scopeWithRels, fcr, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc))

	if errorkit.IsNotNilThenLog(detailedErr) {
		if !detailedErr.Flow {
			w.WriteHeader(http.StatusInternalServerError)
		} else if detailedErr.ErrDescConst == entity.FlowErrBearerAuthzTokenNotFound || detailedErr.ErrDescConst == entity.FlowErrUnauthorizedBearerAuthzToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(adapter.MakeFinishClientRegistrationResponseBody(nil, "", []error{detailedErr}, nil, nil))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(adapter.MakeFinishClientRegistrationResponseBody(nil, "", []error{detailedErr}, nil, nil))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(adapter.MakeFinishClientRegistrationResponseBody(nil, "", nil, result, &fcrr.FinishClientRegistrationShared))
}

func (fcr *FinishClientRegistration) getAuthzToken(bearerAuthzToken string) string {
	if !strings.Contains(bearerAuthzToken, "Bearer ") {
		return ""
	}
	split := strings.Split(bearerAuthzToken, " ")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

func (fcr *FinishClientRegistration) validateAndParseRequest(r *http.Request) (fcrr *entity.FinishClientRegistrationRequest, regexErrMsgs map[string][]string) {
	var callTraceFunc = fmt.Sprintf("%s#*FinishClientRegistration.validateAndParseRequest", callTraceFileFinishRegisterClientController)

	headerRules := make(map[string]uint)
	headerRules["Authorization"] = adapter.RegexBearerAuthzToken
	headerRules["Content-Type"] = adapter.RegexAppliationJson
	headerRules["Accept"] = adapter.RegexAppliationJson
	hpv := restkit.HeaderParamValidation{RegexRules: headerRules, Header: r.Header}

	if regexErrMsgs, ok := hpv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen); !ok {
		return nil, *regexErrMsgs
	}

	bodyData := make([]byte, r.ContentLength)
	r.Body.Read(bodyData)
	defer r.Body.Close()

	var result entity.FinishClientRegistrationRequest
	err := json.Unmarshal(bodyData, &result)
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrJsonUnmarshal, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "FinishClientRegistrationRequest"))
		return nil, nil
	}

	if regexErrMsgs = fcr.validateRequestBody(&result); len(regexErrMsgs) > 0 {
		return nil, regexErrMsgs
	}

	return &result, nil
}

func (fcr *FinishClientRegistration) validateRequestBody(body *entity.FinishClientRegistrationRequest) map[string][]string {
	var allErrMsgs map[string][]string = make(map[string][]string)

	var redirectURIsMsgs []string
	for idx := range body.RedirectURIs {
		if !regexkit.RegexpCompiled[regexkit.RegexURL].Match([]byte(body.RedirectURIs[idx])) {
			redirectURIsMsgs = append(redirectURIsMsgs, adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexURL))
		}
	}
	if len(redirectURIsMsgs) > 0 {
		allErrMsgs["redirect_uris"] = redirectURIsMsgs
	}

	if !regexkit.RegexpCompiled[adapter.RegexTokenEndpointAuthMethod].Match([]byte(body.TokenEndpointAuthMethod)) {
		allErrMsgs["token_endpoint_auth_method"] = []string{adapter.RegexErrDescGen.GenerateDesc(adapter.RegexTokenEndpointAuthMethod)}
	}

	var grantTypesMsgs []string
	for idx := range body.GrantTypes {
		if !regexkit.RegexpCompiled[adapter.RegexGrantTypes].Match([]byte(body.GrantTypes[idx])) {
			grantTypesMsgs = append(grantTypesMsgs, adapter.RegexErrDescGen.GenerateDesc(adapter.RegexGrantTypes))
		}
	}
	if len(grantTypesMsgs) > 0 {
		allErrMsgs["grant_types"] = grantTypesMsgs
	}

	var responseTypesMsgs []string
	for idx := range body.ResponseTypes {
		if !regexkit.RegexpCompiled[adapter.RegexResponseTypes].Match([]byte(body.ResponseTypes[idx])) {
			responseTypesMsgs = append(responseTypesMsgs, adapter.RegexErrDescGen.GenerateDesc(adapter.RegexResponseTypes))
		}
	}
	if len(responseTypesMsgs) > 0 {
		allErrMsgs["response_types"] = responseTypesMsgs
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.ClientName)) {
		allErrMsgs["client_name"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexNotEmpty)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexURL].Match([]byte(body.ClientURI)) {
		allErrMsgs["client_uri"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexURL].Match([]byte(body.LogoURI)) {
		allErrMsgs["logo_uri"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.Scope)) {
		allErrMsgs["scope"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexNotEmpty)}
	}

	var contactsMsgs []string
	for idx := range body.Contacts {
		if !regexkit.RegexpCompiled[regexkit.RegexEmail].Match([]byte(body.Contacts[idx])) {
			contactsMsgs = append(contactsMsgs, adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexEmail))
		}
	}
	if len(contactsMsgs) > 0 {
		allErrMsgs["contacts"] = contactsMsgs
	}

	if !regexkit.RegexpCompiled[regexkit.RegexURL].Match([]byte(body.TosURI)) {
		allErrMsgs["tos_uri"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.SoftwareID)) {
		allErrMsgs["software_id"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexNotEmpty)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.SoftwareVersion)) {
		allErrMsgs["software_version"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexNotEmpty)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexUUIDV4].Match([]byte(body.InitClientID)) {
		allErrMsgs["init_client_id"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexUUIDV4)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.Basepoint)) {
		allErrMsgs["basepoint"] = []string{adapter.RegexErrDescGen.GenerateDesc(adapter.RegexBase64RawURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.ClientPK)) {
		allErrMsgs["client_pk"] = []string{adapter.RegexErrDescGen.GenerateDesc(adapter.RegexBase64RawURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexNotEmpty].Match([]byte(body.ServerPK)) {
		allErrMsgs["server_pk"] = []string{adapter.RegexErrDescGen.GenerateDesc(adapter.RegexBase64RawURL)}
	}

	if !regexkit.RegexpCompiled[regexkit.RegexDateTimeRFC3339].Match([]byte(body.SessionExpiredAt.Format(time.RFC3339Nano))) {
		allErrMsgs["session_expired_at"] = []string{adapter.RegexErrDescGen.GenerateDesc(regexkit.RegexDateTimeRFC3339)}
	}

	return allErrMsgs
}

func (fcr *FinishClientRegistration) SelectClientRegistrationsBy(initClientID string) (*entity.ClientRegistration, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*FinishClientRegistration.SelectClientRegistrationBy", callTraceFileFinishRegisterClientController)

	row, err := fcr.dbo.QueryRow("SELECT * FROM client_registrations WHERE init_client_id = ?", initClientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "client_registrations", "init_client_id"); detailedErr != nil {
		return nil, detailedErr
	}

	var celientRegistration entity.ClientRegistration
	if err := row.Scan(&celientRegistration.ID, &celientRegistration.CreatedAt, &celientRegistration.UpdatedAt, &celientRegistration.SoftDeletedAt, &celientRegistration.InitClientID, &celientRegistration.Basepoint, &celientRegistration.ServerSK, &celientRegistration.ServerPK, &celientRegistration.SessionExpiredAt); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(true, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations", "celientRegistration")
	}
	return &celientRegistration, nil
}

type insertScopeWithRel struct {
	scope       string
	parentScope *uint64
}

func (fcr *FinishClientRegistration) InsertClientWithRel(clientsWithRel *entity.ClientWithRel) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#*FinishClientRegistration.InsertClientWithRel", callTraceFileFinishRegisterClientController)

	stmtInsertClient := fmt.Sprintf("INSERT INTO clients (token_endpoint_auth_method, grant_types, response_types, client_name, client_uri, logo_uri, tos_uri, policy_uri, software_id, software_version, init_client_id, client_id, client_id_issued_at, client_secret, client_secret_expired_at) VALUES %s", sqlkit.GeneratePlaceHolder(15))
	clientInsertResult, err := fcr.dbo.Command(stmtInsertClient, clientsWithRel.TokenEndpointAuthMethod, clientsWithRel.GrantTypes, clientsWithRel.ResponseTypes, clientsWithRel.ClientName, clientsWithRel.ClientURI, clientsWithRel.LogoURI, clientsWithRel.TosURI, clientsWithRel.PolicyURI, clientsWithRel.SoftwareID, clientsWithRel.SoftwareVersion, clientsWithRel.InitClientID, clientsWithRel.ClientID, clientsWithRel.ClientIDIssuedAt, clientsWithRel.ClientSecret, clientsWithRel.ClientSecretExpiredAt)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "clients")
	}
	clientsID, err := clientInsertResult.LastInsertId()
	clientsIDUint64 := uint64(clientsID)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBLastInsertId, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "clients")
	}

	stmtInsertRedirectURIs := fmt.Sprintf("INSERT INTO redirect_uris (uri, clients_id) VALUES %s", sqlkit.GenerateNPlaceHolder(uint16(len(clientsWithRel.RedirectURIs)), 2))
	redirecURIs := entity.FlattenRedirectURIsNonTemplateColumnsValue(clientsWithRel.RedirectURIs, uint64(clientsID))
	if _, err = fcr.dbo.Command(stmtInsertRedirectURIs, redirecURIs...); err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "redirect_uris")
	}

	stmtInsertContacts := fmt.Sprintf("INSERT INTO contacts (contact, clients_id) VALUES %s", sqlkit.GenerateNPlaceHolder(uint16(len(clientsWithRel.Contacts)), 2))
	contacts := entity.FlattenContactsNonTemplateColumnsValue(clientsWithRel.Contacts, uint64(clientsID))
	if _, err = fcr.dbo.Command(stmtInsertContacts, contacts...); err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "contacts")
	}

	var parentChildrenIDs map[string][]uint64 = make(map[string][]uint64)
	for _, v := range clientsWithRel.ScopesWithRel {
		result, err := fcr.dbo.Command("INSERT INTO scopes (scope,permission,clients_id) VALUES (?,?,?)", v.Scope.Scope, v.Scope.Permission, clientsID)
		if err != nil {
			return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBInsert, adapter.DetailedErrDescGen)
		}

		lastInsertedID, err := result.LastInsertId()
		if err != nil {
			return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBLastInsertId, adapter.DetailedErrDescGen)
		}
		scopesID := uint64(lastInsertedID)
		v.Scope.ID = &scopesID
		v.Scope.ClientsID = &clientsIDUint64
		if v.ParentScope != nil {
			parentChildrenIDs[v.ParentScope.Scope] = append(parentChildrenIDs[v.ParentScope.Scope], scopesID)
		}
	}

	var txCmdStmtArgs []*sqlkit.TxCmdStmtArgs
	updateParentScopesIDStmt := "UPDATE scopes SET parent_scopes_id=? WHERE id="
	for _, v := range clientsWithRel.ScopesWithRel {
		if _, ok := parentChildrenIDs[v.Scope.Scope]; ok {
			for _, id := range parentChildrenIDs[v.Scope.Scope] {
				txCmdStmtArgs = append(txCmdStmtArgs, &sqlkit.TxCmdStmtArgs{Statement: updateParentScopesIDStmt, Args: []interface{}{v.Scope.ID, id}})
			}
		}
	}
	if len(txCmdStmtArgs) > 0 {
		_, err = fcr.dbo.TxCommand(txCmdStmtArgs)
		if err != nil {
			return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBTxCommand, adapter.DetailedErrDescGen)
		}
	}

	return nil
}

func (fcr *FinishClientRegistration) DeleteClientRegistration(initClientID string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#*FinishClientRegistration.DeleteClientRegistration", callTraceFileFinishRegisterClientController)

	if _, err := fcr.dbo.Command("DELETE FROM client_registrations WHERE init_client_id = ?", initClientID); err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBDelete, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "client_registrations")
	}
	return nil
}

func NewFinishClientRegistration(db *sql.DB) *FinishClientRegistration {
	return &FinishClientRegistration{dbo: &sqlkit.DBOperation{DB: db}}
}
