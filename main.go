package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/ilhammhdd/go-toolkit/errorkit"
	"github.com/ilhammhdd/go-toolkit/regexkit"
	"github.com/ilhammhdd/go-toolkit/restkit"
	"ilhammhdd.com/oauth2-go-server/adapter"
	"ilhammhdd.com/oauth2-go-server/controller"
	"ilhammhdd.com/oauth2-go-server/entity"
	"ilhammhdd.com/oauth2-go-server/external"
	"ilhammhdd.com/oauth2-go-server/usecase"
)

const callTraceFileMain = "/main.go"

func init() {
	var callTraceFunc string = fmt.Sprintf("%s#init", callTraceFileMain)
	totalRequiredArgs := 1 + len(entity.EnvVarKeys)*2
	if len(os.Args) == totalRequiredArgs {
		entity.EnvVars = make(map[string]entity.EnvVarValue)
		for i := 1; i < totalRequiredArgs; i += 2 {
			entity.EnvVars[os.Args[i]] = entity.EnvVarValue{Value: os.Args[i+1], FromArgs: true}
		}
		external.MariaDB = external.OpenDBConnection(entity.EnvVars[entity.DBSourceNameEnvVar].Value, entity.EnvVars[entity.DBUserEnvVar].Value, entity.EnvVars[entity.DBPasswordEnvVar].Value, entity.EnvVars[entity.DBNameEnvVar].Value)
	} else {
		for _, val := range entity.EnvVarKeys {
			entity.EnvVars[val] = entity.EnvVarValue{Value: os.Getenv(val), FromArgs: false}
			defer func(val string) {
				if err := os.Unsetenv(val); err != nil {
					errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrUnsetEnvVar, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc), val))
				}
			}(val)
		}

		external.MariaDB = external.OpenDBConnection(os.Getenv(entity.DBSourceNameEnvVar), os.Getenv(entity.DBUserEnvVar), os.Getenv(entity.DBPasswordEnvVar), os.Getenv(entity.DBNameEnvVar))
	}
	regexkit.CompileAllRegex(adapter.Regex)
	fmt.Printf("\n%v", entity.EnvVars)
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var callTraceFunc string = fmt.Sprintf("%s#main", callTraceFileMain)
	entity.EphemeralKeyPair = usecase.GenerateKeyPair(entity.EnvVars[entity.PlainKeyPairSeedEnvVar].Value)
	if !entity.EnvVars[entity.PlainKeyPairSeedEnvVar].FromArgs {
		os.Unsetenv(entity.PlainKeyPairSeedEnvVar)
	}

	entity.EphemeralSecretsalt = usecase.GenerateSecretSalt(entity.EnvVars[entity.PlainSecretSaltEnvVar].Value)
	if !entity.EnvVars[entity.PlainSecretSaltEnvVar].FromArgs {
		os.Unsetenv(entity.PlainSecretSaltEnvVar)
	}

	entity.EphemeralCsrfTokenKey = usecase.GenerateCsrfTokenKey(entity.EnvVars[entity.CsrfTokenPlainKeyEnvVar].Value)
	if !entity.EnvVars[entity.CsrfTokenPlainKeyEnvVar].FromArgs {
		os.Unsetenv(entity.CsrfTokenPlainKeyEnvVar)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/register-client", &restkit.MethodRouting{
		GetHandler:  controller.NewInitiateClientRegistration(external.MariaDB),
		PostHandler: controller.NewFinishClientRegistration(external.MariaDB),
	})
	registerCORSPolicy := restkit.CORSHeaderPolicy{
		AccessControlAllowOrigin:      "http://localhost:7575",
		AccessControlAllowCredentials: true,
		AccessControlAllowMethods:     restkit.NewCaseSensitiveStrings(restkit.UpperCase, http.MethodGet, http.MethodOptions, http.MethodPost),
		AccessControlAllowHeaders:     restkit.NewCaseSensitiveStrings(restkit.LowerCase, "Authorization", "X-PINGOTHER"),
	}
	// TODO: strange behavior! When the client just refreshed the client secret and trying to hit this endpoint's GET, the browser doesn't send the preflight request, but the subsequent requests does
	http.Handle("/init/register", &restkit.MethodRouting{
		GetHandler: controller.NewGenerateURLOneTimeToken(external.MariaDB),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet: registerCORSPolicy,
		},
	})
	http.Handle("/register", &restkit.MethodRouting{
		GetHandler:  controller.NewRenderRegisterUser(external.MariaDB, controller.HTMLTemplateExecutorFunc(external.ExecuteHTMLTemplate)),
		PostHandler: controller.NewRegisterUser(external.MariaDB, controller.HTMLTemplateExecutorFunc(external.ExecuteHTMLTemplate)),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet:  registerCORSPolicy,
			http.MethodPost: registerCORSPolicy,
		},
	})
	loginCORSPolicy := restkit.CORSHeaderPolicy{
		AccessControlAllowOrigin:      "http://localhost:7575",
		AccessControlAllowCredentials: true,
		AccessControlAllowMethods:     restkit.NewCaseSensitiveStrings(restkit.UpperCase, http.MethodGet, http.MethodOptions, http.MethodPost),
		AccessControlAllowHeaders:     restkit.NewCaseSensitiveStrings(restkit.LowerCase, "Authorization", "X-PINGOTHER"),
	}
	http.Handle("/init/login", &restkit.MethodRouting{
		GetHandler: controller.NewGenerateURLOneTimeToken(external.MariaDB),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet: loginCORSPolicy,
		},
	})
	// TODO: finish this login feature
	http.Handle("/login", &restkit.MethodRouting{
		GetHandler: controller.NewRenderLoginUser(external.MariaDB, controller.HTMLTemplateExecutorFunc(external.ExecuteHTMLTemplate)),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet:  loginCORSPolicy,
			http.MethodPost: loginCORSPolicy,
		},
	})
	refreshClientSecretCORS := restkit.CORSHeaderPolicy{
		AccessControlAllowOrigin:      "http://localhost:7575",
		AccessControlAllowCredentials: true,
		AccessControlAllowMethods:     restkit.NewCaseSensitiveStrings(restkit.UpperCase, http.MethodGet, http.MethodOptions, http.MethodPost),
		AccessControlAllowHeaders:     restkit.NewCaseSensitiveStrings(restkit.LowerCase, "Authorization", "X-PINGOTHER"),
	}
	http.Handle("/refresh-client-secret", &restkit.MethodRouting{
		GetHandler:  controller.NewInitRefreshClientSecret(external.MariaDB),
		PostHandler: controller.NewFinishRefreshClientSecret(external.MariaDB),
		MethodsCORSHeaderPolicy: &restkit.MethodsCORSHeaderPolicy{
			http.MethodGet:  refreshClientSecretCORS,
			http.MethodPost: refreshClientSecretCORS,
		},
	})

	server := http.Server{Addr: fmt.Sprintf(":%s", entity.EnvVars[entity.HTTPPortEnvVar].Value)}
	defer server.Close()
	err := server.ListenAndServe()
	if err != nil {
		errorkit.IsNotNilThenLog(errorkit.NewDetailedError(false, callTraceFunc, err, entity.ErrListenAndServe, errorkit.ErrDescGeneratorFunc(adapter.GenerateDetailedErrDesc)))
	}
}
