package srv_conf

import (
	"log"

	"github.com/spf13/viper"
)

var (
	srvConf = *viper.New()

	fileName string = "srv.yaml"
	fileType string = "yaml"

	AppDir   string
	DataDir  string
	StaticDir string
	AssetsDir string
	ReportsDir string
)

func init() {
	srvConf.SetConfigName(fileName)
	srvConf.AddConfigPath(".")
	srvConf.SetConfigType(fileType)
}

func WriteConfigFile(app_path string) error {

	AppDir = app_path

	srvConf.SetDefault("server_port", 5443)
	srvConf.SetDefault("gin_mode", "release")
	srvConf.SetDefault("app_dir", app_path)
	srvConf.SetDefault("tls_keysize", 4096)

	err := srvConf.WriteConfigAs(fileName)
	if err != nil {
		log.Println("Error creating", fileName)
		return err
	}

	return nil
}

func Srv_ReadConfig() {
	err := srvConf.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading", fileName)
	}
}

func GetString(key string) string {
	return srvConf.GetString(key)
}

func GetInt(key string) int {
	return srvConf.GetInt(key)
}

func GetInt64(key string) int64 {
	return srvConf.GetInt64(key)
}

func GetBool(key string) bool {
	return srvConf.GetBool(key)
}

func SetVal(key string, val any) {
	srvConf.Set(key, val)
	err := srvConf.WriteConfigAs(fileName)
	if err != nil {
		log.Fatal("Error SetVal", fileName)
	}
}

func IsGinModDebug() bool {
	return GetString("gin_mode") == "debug"
}

func SetPaths() {
	DataDir = "data"
	AssetsDir = "appfiles"
	StaticDir = "static"
	ReportsDir = AssetsDir + "/reports"
}

func TLSKeySize() int {
	return GetInt("tls_keysize")
}