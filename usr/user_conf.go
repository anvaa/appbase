package users

import (
	"log"

	"github.com/spf13/viper"
)

var (
	UsrConf  = viper.New()
	fileName = "usr.yaml"
	fileType = "yaml"
)

func init() {
	UsrConf.SetConfigName(fileName)
	UsrConf.AddConfigPath(".")
	UsrConf.SetConfigType(fileType)
}

func WriteConfigFile(appPath string) {
	UsrConf.SetDefault("app_dir", appPath)
	UsrConf.SetDefault("access_time", "3600*12*1") // 12 hours

	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error creating %s: %v", fileName, err)
	}
}

func ReadConfig() {
	if err := UsrConf.ReadInConfig(); err != nil {
		log.Fatalf("Error reading %s: %v", fileName, err)
	}
}

func GetString(key string) string {
	return UsrConf.GetString(key)
}

func GetInt(key string) int {
	return UsrConf.GetInt(key)
}

func GetInt64(key string) int64 {
	return UsrConf.GetInt64(key)
}

func GetBool(key string) bool {
	return UsrConf.GetBool(key)
}

func SetVal(key string, val any) {
	UsrConf.Set(key, val)
	if err := UsrConf.WriteConfigAs(fileName); err != nil {
		log.Fatalf("Error setting value in %s: %v", fileName, err)
	}
}
