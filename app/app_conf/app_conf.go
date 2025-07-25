package app_conf

import (
	"app/app_models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var (
	appConf         = *viper.New()
	fileName string = "app.yaml"
	fileType string = "yaml"

	CookieName          = strings.ReplaceAll(AppName, " ", "") + "_Auth"
	CookieExpire int    = 24 * 60 * 60 // 24 hours 
	UserKey      string = GetUserKey()

	DBName string
)

func init() {
	appConf.SetConfigName(fileName)
	appConf.AddConfigPath(".")
	appConf.SetConfigType(fileType)
}

func WriteDefaultConfig(appRoot string) {
	// SetDefault sets the default value for the key.
	appConf.SetDefault("app_db", "data/app.db")
	appConf.SetDefault("app_dir", appRoot)
	appConf.SetDefault("print_height", 24)
	appConf.SetDefault("print_width", 160)
	appConf.SetDefault("print_margin", 3)
	appConf.SetDefault("print_font_size", 16)
	appConf.SetDefault("print_txt", 1)
	appConf.SetDefault("session_expire", time.Hour*24*1) // 1 day
	appConf.SetDefault("start_url", "/app/start") // Default start page

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

func PrintConf() app_models.PrinterConf {
	// Get the print size from the config
	conf := app_models.PrinterConf{
		Height:   GetInt("print_height"),
		Width:    GetInt("print_width"),
		Margin:   GetInt("print_margin"),
		FontSize: GetInt("print_font_size"),
		Prnttxt:  GetInt("print_txt"),
	}
	return conf
}

func AppConf() app_models.AppConf {
	// Get the app config
	conf := app_models.AppConf{
		StartPageFocus: GetString("start_page_focus"),
	}
	return conf
}

func SessionExpire() time.Duration {
	// Get the session expire time from the config
	return GetDuration("session_expire")
}

func GetDBName() string {
	// Get app_dir and app_db from the 
	dbname := fmt.Sprintf("%s/%s", GetString("start_url"), GetString("app_db"))
	return dbname
}

func GetUserKey() string {
	// Generate a unique user key if not set
	return fmt.Sprintf("%d", uuid.New().ID())
}