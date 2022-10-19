package entity

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/argon2"
)

const callTraceFileUser = "/entity/user.go"

type User struct {
	Id            uint64     `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	Email         string     `json:"email,omitempty"`
	Password      string     `json:"password,omitempty"`
	ClientsID     uint64     `json:"clients_id,omitempty"`
}

func NewUser(email string, plainPassword *string, clientsID uint64, errDescGen errorkit.ErrDescGenerator) (*User, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#NewUser", callTraceFileUser)
	randSalt, err := GenerateCryptoRand(32)
	if err != nil {
		return nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "users")
	}
	cipherPassword := argon2.Key([]byte(*plainPassword), randSalt, 6, 96*1024, 6, 32)
	defer func() { *plainPassword = "" }()

	return &User{
		Email: email, Password: base64.RawURLEncoding.EncodeToString(cipherPassword), ClientsID: clientsID,
	}, nil
}

type Username struct {
	Id            uint64     `json:"id,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	SoftDeletedAt *time.Time `json:"soft_deleted_at,omitempty"`
	UnqNum        uint16     `json:"unq_num,omitempty"`
	Username      string     `json:"username,omitempty"`
	UsersID       uint64     `json:"users_id,omitempty"`
}

type UserWithRel struct {
	User     *User     `json:"user,omitempty"`
	Username *Username `json:"username,omitempty"`
}
