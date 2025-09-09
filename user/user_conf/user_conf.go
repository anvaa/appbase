package user_conf

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"app/app_conf"
	"server/global"
)

var (
	UsrConf    = viper.New()
	fileName   = "usr.yaml"
	fileType   = "yaml"
	CookieName = strings.ReplaceAll(app_conf.AppInfo().AppName, " ", "") + "_Auth"
	CookieExpire = 24 * 60 * 60 // 24 hours
	UserKey      = global.UuidToString(uuid.New().ID())
)

func init() {
	UsrConf.SetConfigName(fileName)
	UsrConf.AddConfigPath(".")
	UsrConf.SetConfigType(fileType)
}

func WriteConfigFile(appPath string) {
	UsrConf.SetDefault("app_dir", appPath)
	UsrConf.SetDefault("session_expire", 12*time.Hour) // 1/2 day
	UsrConf.SetDefault("login_rate_limit", 60)         // in seconds

	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error creating %s: %v", fileName, err)
	}
}

func ReadConfig() {
	if err := UsrConf.ReadInConfig(); err != nil {
		log.Fatalf("Error reading %s: %v", fileName, err)
	}
}

func GetString(key string) string      { return UsrConf.GetString(key) }
func GetInt(key string) int            { return UsrConf.GetInt(key) }
func GetInt64(key string) int64        { return UsrConf.GetInt64(key) }
func GetBool(key string) bool          { return UsrConf.GetBool(key) }

func SetVal(key string, val any) {
	UsrConf.Set(key, val)
	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error setting value in %s: %v", fileName, err)
	}
}

func SessionExpire() time.Duration {
	return UsrConf.GetDuration("session_expire")
}

func LoginRateLimit() time.Duration {
	return UsrConf.GetDuration("login_rate_limit")
}
