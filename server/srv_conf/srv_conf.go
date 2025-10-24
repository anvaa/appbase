package srv_conf

import (
	"app/app_models"

	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	srvConf = *viper.New()

	fileName string = "srv.yaml"
	fileType string = "yaml"

	AppDir     string
	DataDir    string
	StaticDir  string
	AssetsDir  string
	ReportsDir string
	LogsDir    string
)

func init() {
	srvConf.SetConfigName(fileName)
	srvConf.AddConfigPath(".")
	srvConf.SetConfigType(fileType)
}

func WriteConfigFile(app_path string) error {

	AppDir = app_path

	srvConf.SetDefault("host_ip", "")
	srvConf.SetDefault("host_name", "")
	srvConf.SetDefault("web_server_port", "5443")

	srvConf.SetDefault("gin_mode", "release")
	srvConf.SetDefault("use_tls", true)

	srvConf.SetDefault("app_dir", app_path)
	srvConf.SetDefault("tls_keysize", 4096)
	srvConf.SetDefault("MaxLogSizeMB", 10)

	// Database configuration
	srvConf.SetDefault("app_db", "data/app.db")
	srvConf.SetDefault("db_type", "sqlite") // sqlite, mysql or postgres
	srvConf.SetDefault("db_host", "localhost")
	srvConf.SetDefault("db_port", "3306")
	srvConf.SetDefault("db_user", "user")
	srvConf.SetDefault("db_password", "password")
	srvConf.SetDefault("db_name", "dbname")

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

func GetAny(key string) any {
	return srvConf.Get(key)
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

func MaxLogSizeMB() int64 {
	return GetInt64("MaxLogSizeMB")
}

func SetPaths() {
	DataDir = "data"
	AssetsDir = "appfiles"
	StaticDir = "static"
	ReportsDir = AssetsDir + "/reports"
	LogsDir = AssetsDir + "/logs"
}

func TLSKeySize() int {
	return GetInt("tls_keysize")
}

func DBType() string {
	return GetString("db_type")
}

func AppDBPath() string {
	return GetString("app_db")
}

func GetDbConfig() app_models.DbConfig {
	return app_models.DbConfig{
		Type:     DBType(),
		Path:     AppDBPath(),
		Host:     GetString("db_host"),
		Port:     GetString("db_port"),
		User:     GetString("db_user"),
		Password: GetString("db_password"),
		Name:     GetString("db_name"),
	}
}

func SetDBConfig(config app_models.DbConfig) error {
	srvConf.Set("db_type", config.Type)
	srvConf.Set("app_db", config.Path)
	srvConf.Set("db_host", config.Host)
	srvConf.Set("db_port", config.Port)
	srvConf.Set("db_user", config.User)
	srvConf.Set("db_password", config.Password)
	srvConf.Set("db_name", config.Name)
	err := srvConf.WriteConfigAs(fileName)
	if err != nil {
		log.Fatal("Error SetDBConfig", fileName)
		return err
	}
	return nil
}

func UseTLS() bool {
	return GetBool("use_tls")
}

func GetDbInfo() any {
	ret := map[string]string{
		"dbtype": GetString("db_type"),
		"dbpath": GetString("app_dir") + "/" + GetString("app_db"),
		"dbhost": GetString("db_host"),
		"dbport": GetString("db_port"),
		"dbname": GetString("db_name"),
	}
	return ret
}

func GetHostInfo() any {
	ret := map[string]string{
		"ip":   GetHostIP(0)[0],
		"name": GetHostName(),
		"port": GetWebServerPort(),
	}
	return ret
}

func GetWebServerPort() string {
	return GetString("web_server_port")
}

func GetHostName() string {
	// split on '.' and return [0]
	return strings.Split(GetString("host_name"), ".")[0]
}

func GetHostIP(index int) []string {
	_ips := GetAny("host_ip").([]string)
	if len(_ips) == 0 {
		return []string{"noip"}
	}
	if index < 0 || index >= len(_ips) {
		return []string{_ips[0]}
	}
	ret := _ips[index]

	return []string{ret}
}
