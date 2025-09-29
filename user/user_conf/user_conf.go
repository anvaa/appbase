package user_conf

import (
	"log"
	"time"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"server/global"
)

var (
	UsrConf    = viper.New()
	fileName   = "usr.yaml"
	fileType   = "yaml"
	CookieName = GenCookieName()
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
	UsrConf.SetDefault("session_expire", 12) // 12 hours
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
	h := GetInt("session_expire")
	return time.Duration(time.Now().Add(time.Duration(h) * time.Hour).Unix())
}

func LoginRateLimit() time.Duration {
	return UsrConf.GetDuration("login_rate_limit")
}

func GenCookieName() string {
	// Make cookie name URL friendly
	// Replace spaces with underscores and convert to lowercase
	appName := fmt.Sprintf("auth_appbase_%v", time.Now().UnixNano())

	return appName
}