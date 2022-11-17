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

type RenderLoginUser struct {
	dbo                  DBOperator
	htmlTemplateExecutor HTMLTemplateExecutor
}

func (rlu RenderLoginUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	theUsecase := usecase.UrlOneTimeTokenUsecase{
		ErrDescGen: adapter.DetailedErrDescGen, DBO: UrlOneTimeToken{rlu.dbo},
	}

	verified, detailedErr := theUsecase.VerifySignature(r.URL.Query().Get("client_id"), r.URL.Query().Get("one_time_token"), r.URL.Query().Get("signature"))
	if errorkit.IsNotNilThenLog(detailedErr) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}

	if !verified {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "invalid token", nil))
		return
	}

	csrfToken, csrfTokenHmac, detailedErr := usecase.GenerateCsrfTokenAndHmac(adapter.DetailedErrDescGen)
	if errorkit.IsNotNilThenLog(detailedErr) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(adapter.MakeResponseTmplErrResponse(nil, "", []error{detailedErr}))
		return
	}

	tmplData := registerAndLoginTmplData{CsrfToken: csrfToken, ClientID: r.URL.Query().Get("client_id")}

	http.SetCookie(w, &http.Cookie{Name: "csrf-token-hmac-login", Value: csrfTokenHmac, Path: "/login", MaxAge: 300, HttpOnly: true, Secure: true})
	w.WriteHeader(http.StatusOK)
	rlu.htmlTemplateExecutor.ExecuteTemplate("login_tmpl", "static/login.html", tmplData, w)
}

func NewRenderLoginUser(db *sql.DB, htmlTemplateExecutor HTMLTemplateExecutor) RenderLoginUser {
	return RenderLoginUser{dbo: sqlkit.DBOperation{DB: db}, htmlTemplateExecutor: htmlTemplateExecutor}
}
