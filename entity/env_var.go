package entity

const (
	PlainKeyPairSeedEnvVar                           = "PLAIN_KEY_PAIR_SEED"
	PlainSecretSaltEnvVar                            = "PLAIN_SECRET_SALT"
	DBSourceNameEnvVar                               = "DB_SOURCE_NAME"
	DBUserEnvVar                                     = "DB_USER"
	DBPasswordEnvVar                                 = "DB_PASSWORD"
	DBNameEnvVar                                     = "DB_NAME"
	HttpPortEnvVar                                   = "HTTP_PORT"
	OverTLS                                          = "OVER_TLS"
	ClientRegistrationSessionExpirationSecondsEnvVar = "CLIENT_REGISTRATION_SESSION_EXPIRATION_SECONDS"
	CsrfTokenPlainKeyEnvVar                          = "CSRF_TOKEN_PLAIN_KEY"
)

var EnvVarKeys []string = []string{
	DBSourceNameEnvVar,
	DBUserEnvVar,
	DBPasswordEnvVar,
	DBNameEnvVar,
	PlainKeyPairSeedEnvVar,
	PlainSecretSaltEnvVar,
	HttpPortEnvVar,
	OverTLS,
	ClientRegistrationSessionExpirationSecondsEnvVar,
	CsrfTokenPlainKeyEnvVar,
}

type EnvVarValue struct {
	Value    string
	FromArgs bool
}

var EnvVars map[string]EnvVarValue
