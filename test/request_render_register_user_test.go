package test

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/blake2b"
)

func TestPrepareRequest(t *testing.T) {
	decodedRaw, _ := base64.RawURLEncoding.DecodeString("qGMS1XRmfLQZDBuWTRDvhvw3umtEVrcAzAIcvvQjj0k")
	checksum := blake2b.Sum256(decodedRaw)
	checksumEncoded := base64.RawURLEncoding.EncodeToString(checksum[:])
	t.Logf("hashEncoded: %s", checksumEncoded)
}
