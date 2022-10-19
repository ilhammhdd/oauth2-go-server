package entity

import (
	"crypto/rand"
	"errors"
)

const StackTraceSize = 1024 * 8

var DBTableTemplateColumns []string = []string{
	"id",
	"created_at",
	"updated_at",
	"soft_deleted_at",
}

type DBTable interface {
	Columns() []string
}

func GenerateCryptoRand(randLen int) ([]byte, error) {
	randBuff := make([]byte, randLen)
	randN, randErr := rand.Reader.Read(randBuff)
	defer func() {
		randBuff = make([]byte, randLen)
	}()
	if randN != randLen {
		return nil, errors.New("random n != randLen")
	}
	return randBuff, randErr
}
