package entity

import "time"

type FinishRefreshClientSecretRequest struct {
	InitClientID     string     `json:"init_client_id"`
	Basepoint        string     `json:"basepoint"`
	ClientPk         string     `json:"client_pk"`
	ServerPk         string     `json:"server_pk"`
	SessionExpiredAt *time.Time `json:"session_expired_at"`
}

type FinishRefreshClientSecretResponse struct {
	ClientSecretExpiredAt time.Time `json:"client_secret_expired_at"`
	ResponseBodyTemplate
}
