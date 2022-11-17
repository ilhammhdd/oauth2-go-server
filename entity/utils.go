package entity

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	mathRand "math/rand"
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

func GenerateRandID() string {
	var mathBuf []byte
	mathBuf = binary.LittleEndian.AppendUint64(mathBuf, uint64(mathRand.Int63n(281_474_976_710_656)))

	randID := []rune(base64.RawURLEncoding.EncodeToString(mathBuf))
	randID = randID[:8]

	for i := range randID {
		if randID[i] == '-' || randID[i] == '_' {
			randID[i] = rune(mathRand.Int31n(26) + 65)
		}
	}

	return string(randID)
}
