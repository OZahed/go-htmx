/*
In this use case we only need Int, bool, string values for our config, if you need more data types for
your configs consider replacing MakeConfigsFromEnv implementation with `viper` initialization and replace the EnvGetter
with Viper object,

Since Golang type system won't allow

If you do not need config default values, consider using `GET[T any](name string) T` function
*/
package config

import (
	"time"

	"github.com/OZahed/bob/envs"
)

type AppConfig struct {
	AppName            string        `env:"NAME,default=HTMX"`
	LayoutsRootDir     string        `env:"LAYOUT_TEMP_DIR,default=./templates"`
	LayoutRootTmpName  string        `env:"LAYOUT_TEMP_ROOT_NAME,default=Layout"`
	PartialRootDirs    string        `env:"PARTIAL_TEMP_DIR,default=./templates/partials"`
	StaticFilesDir     string        `env:"STATIC_FILES_DIR,default=./public"`
	StaticRoutesPrefix string        `env:"STATIC_ROUTES_PREFIX,default=/public/"`
	Mode               string        `env:"MODE,default=DEV"`
	DebuggerBaseName   string        `env:"DEBUGGER_NAME,default=main"`
	CertFile           string        `env:"CERT_FILE,default=server.crt"`
	KeyFile            string        `env:"KEY_FILE,default=server.key"`
	ShutdownDuration   time.Duration `env:"SHUT_DOWN_DURATION,default=5s"`
	TimeOut            time.Duration `env:"TIMEOUT,default=15s"`
	Port               int           `env:"PORT,default=8080"`
}

// NewAppConfig makes the config based onf ENVs
// usually it is best practice to keep your application environment variables like this: APPNAME_ENV_NAME
// prefix is used to take care of APPNAME_ part
func NewAppConfig(Prefix string) (*AppConfig, error) {
	cfg := AppConfig{}

	err := envs.NewParser(envs.DefaultKeyFunc, envs.DefaultGetFunc).ParseStruct(&cfg, Prefix)
	if err != nil {
		return nil, err
	}

	return &cfg, nil

}
