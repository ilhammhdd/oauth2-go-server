package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileRenderRegisterUserController = "/controller/render_register_user_controller.go"

type renderRegisterTmplData struct {
	CsrfToken string
	ClientID  string
}

type RenderRegisterUser struct {
	dbo                  DBOperator
	htmlTemplateExecutor HTMLTemplateExecutor
}

func (rru *RenderRegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var callTraceFunc = fmt.Sprintf("%s#RenderRegisterUser.ServeHTTP", callTraceFileRegisterUserController)

	renderRegisterUserUsecase := usecase.RenderRegisterUser{
		ClientID:     r.URL.Query().Get("client_id"),
		OneTimeToken: r.URL.Query().Get("one_time_token"),
		ReqSignature: r.URL.Query().Get("signature"),
		DBO:          rru,
		ErrDescGen:   errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc),
	}

	verified, detailedErr := renderRegisterUserUsecase.VerifySignature()
	if !verified || errorkit.IsNotNilThenLog(detailedErr) {
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	csrfToken, hmac, detailedErr := renderRegisterUserUsecase.GenerateCsrfTokenAndHmac()
	if errorkit.IsNotNilThenLog(detailedErr) && detailedErr.Flow {
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "csrf-token-hmac", Value: hmac, Path: "/register", MaxAge: 300, HttpOnly: true, Secure: true,
	})
	w.WriteHeader(http.StatusOK)
	rru.htmlTemplateExecutor.ExecuteTemplate("register_tmpl", "static/register.html", renderRegisterTmplData{csrfToken, r.URL.Query().Get("client_id")}, w)
}

func (rru *RenderRegisterUser) DeleteURLOneTimeToken(oneTimeToken string, signature string) *errorkit.DetailedError {
	var callTraceFunc = fmt.Sprintf("%s#*RenderRegisterUser.DeleteURLOneTimeToken", callTraceFileRenderRegisterUserController)

	_, err := rru.dbo.Command("DELETE FROM url_one_time_tokens WHERE one_time_token = ? AND signature = ?", oneTimeToken, signature)
	if err != nil {
		return errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBDelete, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "url_one_time_tokens")
	}
	return nil
}

func (rru *RenderRegisterUser) SelectURLOneTimeToken(clientID string, oneTimeToken string, signature string) (*entity.URLOneTimeToken, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#*RenderRegisterUser.SelectURLOneTimeToken", callTraceFileRenderRegisterUserController)

	row, err := rru.dbo.QueryRow("SELECT id FROM clients WHERE client_id = ?", clientID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "clients", "client_id"); detailedErr != nil {
		return nil, detailedErr
	}

	var clientsID uint64
	if err = row.Scan(&clientsID); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "COUNT(clients.id)", "id")
	}

	row, err = rru.dbo.QueryRow("SELECT * FROM url_one_time_tokens WHERE one_time_token = ? AND signature = ? AND clients_id = ?", oneTimeToken, signature, clientsID)
	if detailedErr := handleSelectTableErr(err, callTraceFunc, "url_one_time_tokens", "one_time_token", "signature", "clients_id"); detailedErr != nil {
		return nil, detailedErr
	}

	var urlOneTimeToken entity.URLOneTimeToken
	if err = row.Scan(&urlOneTimeToken.ID, &urlOneTimeToken.CreatedAt, &urlOneTimeToken.UpdatedAt, &urlOneTimeToken.SoftDeletedAt, &urlOneTimeToken.Pk, &urlOneTimeToken.Sk, &urlOneTimeToken.OneTimeToken, &urlOneTimeToken.Signature, &urlOneTimeToken.URL, &urlOneTimeToken.ClientsID); err != nil && err != sql.ErrNoRows {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrDBScan, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), "url_one_time_tokens.*", "urlOneTimeToken")
	}
	return &urlOneTimeToken, nil
}

func NewRenderRegisterUser(db *sql.DB, htmlTemplateExecutor HTMLTemplateExecutor) *RenderRegisterUser {
	return &RenderRegisterUser{
		dbo: &sqlkit.DBOperation{DB: db}, htmlTemplateExecutor: htmlTemplateExecutor,
	}
}
