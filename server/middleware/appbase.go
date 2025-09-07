package middleware

import (
	"app/app_conf"
	"app/app_models"
	"user/user_conf"

	"github.com/gin-gonic/gin"
)

func AppBase(c *gin.Context) app_models.Appbase {
	return app_models.Appbase{
		Title:   app_conf.AppInfo().AppName,
		Logos:   app_conf.AppLogos(),
		User:    c.Keys[user_conf.UserKey],
		Appinfo: app_conf.AppInfo(),
	}
}
