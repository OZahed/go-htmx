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
)

type AppConfig struct {
	Getter             *Getter
	AppName            string
	LayoutsRootDir     string
	LayoutRootTmpName  string
	PartialRootDirs    string
	StaticFilesDir     string
	StaticRoutesPrefix string
	Mode               string
	DebuggerBaseName   string
	CertFile           string
	KeyFile            string
	ShutdownDuration   time.Duration
	TimeOut            time.Duration
	Port               int
}

// NewAppConfig makes the config based onf ENVs
// usually it is best practice to keep your application environment variables like this: APPNAME_ENV_NAME
// prefix is used to take care of APPNAME_ part
func NewAppConfig(Prefix string) *AppConfig {
	keyProvider := makeKeyProviderPrefix(Prefix)
	getter := &Getter{
		keyProvider: keyProvider,
	}

	cfg := new(AppConfig)
	cfg.Port = GetDefault(keyProvider("PORT"), 8080)
	cfg.Mode = GetDefault(keyProvider("MODE"), "DEV")
	cfg.AppName = GetDefault(keyProvider("NAME"), Prefix)
	cfg.TimeOut = GetDefault(keyProvider("TIMEOUT"), time.Second*15)
	cfg.CertFile = GetDefault(keyProvider("CERT_FILE"), "server.crt")
	cfg.KeyFile = GetDefault(keyProvider("KEY_FILE"), "server.key")
	cfg.LayoutsRootDir = GetDefault(keyProvider("LAYOUT_TEMP_DIR"), "./templates")
	cfg.StaticFilesDir = GetDefault(keyProvider("STATIC_FILES_DIR"), "./public")
	cfg.PartialRootDirs = GetDefault(keyProvider("PARTIAL_TEMP_DIR"), "./templates/partials")
	cfg.DebuggerBaseName = GetDefault(keyProvider("DEBUGGER_NAME"), "main")
	cfg.ShutdownDuration = GetDefault(keyProvider("SHUT_DOWN_DURATION"), time.Second*5)
	cfg.LayoutRootTmpName = GetDefault(keyProvider("LAYOUT_TEMP_ROOT_NAME"), "Layout")
	cfg.StaticRoutesPrefix = GetDefault(keyProvider("STATIC_ROUTES_PREFIX"), "/public/")

	// make sure this part is assigned
	cfg.Getter = getter

	return cfg
}
