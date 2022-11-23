package entity

import "time"

type AuthzCode struct {
	AuthzCodeID string    `json:"authz_code_id"`
	Status      string    `json:"status"`
	StatusAt    time.Time `json:"status_at"`
	TableTemplateCols
}

type AuthzCodeWithRel struct {
	AuthzCode             *AuthzCode             `json:"authz_code"`
	Scopes                []*ScopeWithRel        `json:"scopes"`
	AccessXorRefreshToken *AccessXorRefreshToken `json:"access_and_refresh_tokens"`
	TableTemplateCols
}
