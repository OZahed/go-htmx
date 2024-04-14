/*
In this use case we only need Int, bool, string values for our config, if you need more data types for
your configs consider replacing MakeConfigsFromEnv implementation with `viper` initialization and replace the EnvGetter
with Viper object,

Since Golang type system won't allow
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
	cfg.Port = getter.GetInt("PORT", 8080)
	cfg.Mode = getter.GetString("MODE", "DEV")
	cfg.AppName = getter.GetString("NAME", Prefix)
	cfg.TimeOut = getter.GetDuration("TIMEOUT", time.Second*15)
	cfg.StaticFilesDir = getter.GetString("STATIC_FILES_DIR", "./public")
	cfg.StaticRoutesPrefix = getter.GetString("STATIC_ROUTES_PREFIX", "/public/")
	cfg.LayoutsRootDir = getter.GetString("LAYOUT_TEMP_DIR", "./templates")
	cfg.PartialRootDirs = getter.GetString("PARTIAL_TEMP_DIR", "./templates/partials")
	cfg.DebuggerBaseName = getter.GetString("DEBUGGER_NAME", "main")
	cfg.LayoutRootTmpName = getter.GetString("LAYOUT_TEMP_ROOT_NAME", "Layout")
	cfg.ShutdownDuration = getter.GetDuration("SHUT_DOWN_DURATION", time.Second*5)
	cfg.CertFile = getter.GetString("CERT_FILE", "server.crt")
	cfg.KeyFile = getter.GetString("KEY_FILE", "server.key")

	// make sure this part is assigned
	cfg.Getter = getter

	return cfg
}
