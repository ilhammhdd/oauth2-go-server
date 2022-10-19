package test

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/curve25519"
	"ilhammhdd.com/oauth2-go-server/entity"
)

func TestFinishClientRegistration(t *testing.T) {
	clientSk, err := entity.GenerateCryptoRand(32)
	if err != nil {
		t.Fatalf("\n%s", err.Error())
	}

	basepoint, err := base64.RawURLEncoding.DecodeString("LbYB9Z3GumihliLE07taoEO3z0gCXop4i6QHI_yEF6o")
	if err != nil {
		t.Fatalf("\n%s", err.Error())
	}

	clientPk, err := curve25519.X25519(clientSk, basepoint)
	if err != nil {
		t.Fatalf("\n%s", err.Error())
	}
	t.Logf("\nclientPk: %s", base64.RawURLEncoding.EncodeToString(clientPk))

	serverPk, err := base64.RawURLEncoding.DecodeString("BrEOf4hrlZSXxWepffL1grElnup4gBraYTsH69exbw4")
	if err != nil {
		t.Fatalf("\n%s", err.Error())
	}

	sharedClientSecret, err := curve25519.X25519(clientSk, serverPk)
	if err != nil {
		t.Fatalf("\n%s", err.Error())
	}
	t.Logf("\nsharedClientSecret: %s", base64.RawURLEncoding.EncodeToString(sharedClientSecret))
}

func TestGenRandPass(t *testing.T) {
	r, _ := entity.GenerateCryptoRand(12)
	t.Logf("%s", base64.RawURLEncoding.EncodeToString(r))
}
