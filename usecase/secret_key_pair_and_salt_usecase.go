package usecase

import (
	"crypto/ed25519"
	"encoding/base64"

	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

func GenerateKeyPair(seedEnvVar string) *entity.KeyPair {
	seedBlake2b256Sum := blake2b.Sum256([]byte(seedEnvVar))
	privateKey := ed25519.NewKeyFromSeed(seedBlake2b256Sum[:])
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])
	defer func() {
		seedBlake2b256Sum = [32]byte{}
		privateKey = make([]byte, ed25519.PrivateKeySize)
		publicKey = make([]byte, ed25519.PublicKeySize)
	}()

	privateKeyBase64 := base64.RawURLEncoding.EncodeToString(privateKey)
	publicKeyBase64 := base64.RawURLEncoding.EncodeToString(publicKey)

	return &entity.KeyPair{
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
		PrivateKeyEncoded: privateKeyBase64,
		PublicKeyEncoded:  publicKeyBase64,
	}
}

func GenerateSecretSalt(plainSaltEnvVar string) *entity.SecretSalt {
	secretSaltBlake2b256Sum := blake2b.Sum256([]byte(plainSaltEnvVar))
	secretSaltBase64 := base64.RawURLEncoding.EncodeToString(secretSaltBlake2b256Sum[:])

	return &entity.SecretSalt{Checksum: secretSaltBlake2b256Sum[:], ChecksumEncoded: secretSaltBase64}
}

func GenerateCsrfTokenKey(plainKey string) *entity.CsrfTokenKey {
	csrfTokenKeyBlake2b256Sum := blake2b.Sum256([]byte(plainKey))
	csrfTokenKeyBase64 := entity.CsrfTokenKey(base64.RawURLEncoding.EncodeToString(csrfTokenKeyBlake2b256Sum[:]))

	return &csrfTokenKeyBase64
}
