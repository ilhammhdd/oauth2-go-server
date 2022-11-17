package controller

import (
	"database/sql"
	"net/http"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"github.com/ilhammhdd/go-toolkit/sqlkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileRenderRegisterUserController = "/controller/render_register_user_controller.go"

type RenderRegisterUser struct {
	dbo                  DBOperator
	htmlTemplateExecutor HTMLTemplateExecutor
}

func (rru *RenderRegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// var callTraceFunc = fmt.Sprintf("%s#RenderRegisterUser.ServeHTTP", callTraceFileRegisterUserController)

	uqv := restkit.URLQueryValidation{RegexRules: map[string]uint{
		"client_id":      adapter.RegexRandomID,
		"one_time_token": adapter.RegexRandomID,
		"signature":      regexkit.RegexNotEmpty,
	}, Values: r.URL.Query()}

	regexNotMatchMesgs, ok := uqv.Validate(adapter.DetailedErrDescGen, adapter.RegexErrDescGen)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(regexNotMatchMesgs, "", nil))
		return
	}

	urlOneTimeTokenController := UrlOneTimeToken{rru.dbo}
	urlOneTimeTokenUsecase := usecase.UrlOneTimeTokenUsecase{
		ErrDescGen: adapter.DetailedErrDescGen, DBO: urlOneTimeTokenController,
	}

	verified, detailedErr := urlOneTimeTokenUsecase.VerifySignature(r.URL.Query().Get("client_id"), r.URL.Query().Get("one_time_token"), r.URL.Query().Get("signature"))
	if errorkit.IsNotNilThenLog(detailedErr) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}
	if !verified {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "unverified one time token", nil))
		return
	}

	csrfToken, csrfTokenHmac, detailedErr := usecase.GenerateCsrfTokenAndHmac(adapter.DetailedErrDescGen)
	if errorkit.IsNotNilThenLog(detailedErr) && detailedErr.Flow {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "csrf-token-hmac-register", Value: csrfTokenHmac, Path: "/register", MaxAge: 300, HttpOnly: true, Secure: true,
	})
	w.WriteHeader(http.StatusOK)
	rru.htmlTemplateExecutor.ExecuteTemplate("register_tmpl", "static/register.html", registerAndLoginTmplData{csrfToken, r.URL.Query().Get("client_id")}, w)
}

func NewRenderRegisterUser(db *sql.DB, htmlTemplateExecutor HTMLTemplateExecutor) *RenderRegisterUser {
	return &RenderRegisterUser{
		dbo: &sqlkit.DBOperation{DB: db}, htmlTemplateExecutor: htmlTemplateExecutor,
	}
}
