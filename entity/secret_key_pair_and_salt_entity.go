package entity

const (
	PrivateKeyGeneral = "PRIVATE_KEY_GENERAL"
	PublicKeyGeneral  = "PUBLIC_KEY_GENERAL"
)

type KeyPair struct {
	PrivateKey        []byte
	PublicKey         []byte
	PrivateKeyEncoded string
	PublicKeyEncoded  string
}

var EphemeralKeyPair *KeyPair

type SecretSalt struct {
	Plain           string
	Checksum        []byte
	ChecksumEncoded string
}

var EphemeralSecretsalt *SecretSalt

type CsrfTokenKey string

var EphemeralCsrfTokenKey *CsrfTokenKey
