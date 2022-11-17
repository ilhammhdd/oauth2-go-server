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

func (u *User) VerifyPassword(inPassword string, passwordParams UserPasswordParams) (bool, *errorkit.DetailedError) {
	return false, nil
}

func NewUserWithPasswordParams(email string, plainPassword *string, clientsID uint64, errDescGen errorkit.ErrDescGenerator) (*User, *UserPasswordParams, *errorkit.DetailedError) {
	var callTraceFunc = fmt.Sprintf("%s#NewUserWithPasswordParams", callTraceFileUser)
	randSalt, err := GenerateCryptoRand(32)
	if err != nil {
		return nil, nil, errorkit.NewDetailedError(false, callTraceFunc, err, ErrGenerateCryptoRand, errDescGen, "users")
	}
	var time uint32 = 6
	var memory uint32 = 9 * 1024
	var threads uint8 = 6
	var keyLen uint32 = 32
	var cipherPassword = argon2.Key([]byte(*plainPassword), randSalt, time, memory, threads, keyLen)
	defer func() { *plainPassword = "" }()

	return &User{
		Email: email, Password: base64.RawURLEncoding.EncodeToString(cipherPassword), ClientsID: clientsID,
	}, &UserPasswordParams{RandSalt: base64.RawURLEncoding.EncodeToString(randSalt), Time: time, Memory: memory, Threads: threads, KeyLen: keyLen}, nil
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
	User               *User               `json:"user,omitempty"`
	Username           *Username           `json:"username,omitempty"`
	UserPasswordParams *UserPasswordParams `json:"user_password_params,omitempty"`
}
