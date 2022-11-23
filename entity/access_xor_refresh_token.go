package entity

import "time"

type AccessXorRefreshToken struct {
	Token_type   string    `json:"token_type"`
	Token_id     string    `json:"token_id"`
	Status       string    `json:"status"`
	StatusAt     time.Time `json:"status_at"`
	AuthzCodesID uint64    `json:"authz_codes_id"`
	TableTemplateCols
}
