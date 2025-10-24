package user_conf

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	UsrConf      = viper.New()
	fileName     = "usr.yaml"
	fileType     = "yaml"
	CookieName   = GenCookieName()
	CookieExpire = 24 * 60 * 60 // 24 hours
	UserKey      = "auth_user"  // Fixed key for storing user in context
)

func init() {
	UsrConf.SetConfigName(fileName)
	UsrConf.AddConfigPath(".")
	UsrConf.SetConfigType(fileType)
}

func WriteConfigFile(appPath string) {
	UsrConf.SetDefault("app_dir", appPath)
	UsrConf.SetDefault("session_expire", 12)   // 12 hours
	UsrConf.SetDefault("login_rate_limit", 60) // in seconds

	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error creating %s: %v", fileName, err)
	}
}

func ReadConfig() {
	if err := UsrConf.ReadInConfig(); err != nil {
		log.Fatalf("Error reading %s: %v", fileName, err)
	}
}

func GetString(key string) string { return UsrConf.GetString(key) }
func GetInt(key string) int       { return UsrConf.GetInt(key) }
func GetInt64(key string) int64   { return UsrConf.GetInt64(key) }
func GetBool(key string) bool     { return UsrConf.GetBool(key) }

func SetVal(key string, val any) {
	UsrConf.Set(key, val)
	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error setting value in %s: %v", fileName, err)
	}
}

func SessionExpire() time.Duration {
	h := GetInt("session_expire")
	return time.Duration(h) * time.Hour
}

func LoginRateLimit() time.Duration {
	return UsrConf.GetDuration("login_rate_limit")
}

func GenCookieName() string {
	// Static cookie name that doesn't change between app restarts
	return "_auth_token"
}
