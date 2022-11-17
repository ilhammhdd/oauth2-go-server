package entity

type RegisterUserRequest struct {
	EmailAddress    string `json:"email_address"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
