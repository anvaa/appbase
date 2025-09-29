package app_conf

import (
	
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	appConf         = *viper.New()
	fileName string = "app.yaml"
	fileType string = "yaml"
)

func init() {
	appConf.SetConfigName(fileName)
	appConf.AddConfigPath(".")
	appConf.SetConfigType(fileType)
}

func WriteDefaultConfig(appRoot string) {
	// SetDefault sets the default value for the key.
	appConf.SetDefault("base_url", "/app") // Default start page

	appConf.SetDefault("do_index", false)

	appConf.SetDefault("company_name", "CompanyName")
	appConf.SetDefault("company_url", "https://www.appbase.com")
	appConf.SetDefault("app_name", "AppBase")
	appConf.SetDefault("app_name_long", "AppBase Long Name")
	appConf.SetDefault("app_version", "")
	appConf.SetDefault("logo_small", "logo_small.png")
	appConf.SetDefault("logo_large", "logo_large.png")

	err := appConf.WriteConfigAs(fileName)
	if err != nil {
		log.Fatal("Error writing", fileName)
	}
}

func ReadConfig() {
	err := appConf.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading", fileName)
	}
}

func GetString(key string) string {
	return appConf.GetString(key)
}

func GetInt(key string) int {
	return appConf.GetInt(key)
}

func GetInt64(key string) int64 {
	return appConf.GetInt64(key)
}

func GetTime(key string) time.Time {
	return appConf.GetTime(key)
}

func GetBool(key string) bool {
	return appConf.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return appConf.GetDuration(key)
}

func SetVal(key string, val any) {
	appConf.Set(key, val)
	err := appConf.WriteConfigAs("app.yaml")
	if err != nil {
		log.Fatal("Error SetVal", fileName)
	}
}

func StatusDefault() int {
	return GetInt("status_default")
}

func GetSubDefaults() []int {
	var defs []int
	defs = append(defs, GetInt("sub0_default"))
	defs = append(defs, GetInt("sub1_default"))
	defs = append(defs, GetInt("sub2_default"))
	defs = append(defs, GetInt("sub3_default"))
	defs = append(defs, GetInt("sub4_default"))
	defs = append(defs, GetInt("sub5_default"))
	return defs
}

func LogoSmall() string {
	return GetString("logo_small")
}

func LogoLarge() string {
	return GetString("logo_large")
}

func AppLogos() []string {
	return []string{
		LogoSmall(),
		LogoLarge(),
	}
}

func Doindex() bool {
	return GetBool("do_index")
}

func BaseURL() string {
	return GetString("base_url")
}
