package adapter

type LoginUserByEmailRequest struct {
	EmailAddressUsername string `json:"email_address_username"`
	Password             string `json:"password"`
}
