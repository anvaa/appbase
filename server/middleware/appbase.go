package middleware

import (
	"app/app_conf"
	"app/app_models"
	"fmt"

	"server/srv_conf"

	"user/user_conf"

	"github.com/gin-gonic/gin"
)

func AppBase(c *gin.Context) app_models.Appbase {
	return app_models.Appbase{
		Title:   app_conf.AppInfo().AppName,
		DbInfo:  dbinfo(),
		Doindex: app_conf.Doindex(),
		Logos:   app_conf.AppLogos(),
		User:    c.Keys[user_conf.UserKey],
		Appinfo: app_conf.AppInfo(),
	}
}

func dbinfo() string {
	
	_dbinfo := srv_conf.GetDbInfo().(map[string]string)
	_srvinfo := srv_conf.GetHostInfo().(map[string]string)

	_host := fmt.Sprintf("%s %s", _srvinfo["name"], srv_conf.GetHostIP(0))
	_dbsrv := _srvinfo["db_host"]
	_db := fmt.Sprintf("%s %s@%s:%s", _dbsrv, _dbinfo["dbname"], _dbinfo["dbhost"], _dbinfo["dbport"])

	switch _dbinfo["dbtype"] {
	case "sqlite":
		return fmt.Sprintf("SQLITE: %s %s", _host, _dbinfo["dbpath"])
	case "mysql":
		return "MySQL: " + _db
	case "postgres":
		return "Postgres: " + _db
	}
	return "unknown data source"
}
