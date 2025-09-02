package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"app/app_models"
	"user/user_conf"

	"net/http"

	"github.com/gin-gonic/gin"
)

var appinfo = app_conf.AppInfo()

func appbase(c *gin.Context) app_models.Appbase {
	return app_models.Appbase{
		Title:   app_conf.AppInfo().AppName,
		Logos:   app_conf.AppLogos(),
		User:    c.Keys[user_conf.UserKey],
		Appinfo: appinfo,
	}
}

func Start(c *gin.Context) {
	menu := app_db.Get_MenuTitles()
	c.HTML(http.StatusOK, "start.html", gin.H{
		"appbase": appbase(c),
		"title":   app_conf.AppInfo().AppName,
		"js":      "start.js",
		"menu0":   menu[0],
	})
}

func ToolsTitles(c *gin.Context) {
	menu := app_db.Get_MenuTitles()
	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"appbase": appbase(c),
		"title":   app_conf.AppInfo().AppName + " - Titles",
		"js":      "tools.js",
		"menu":    menu,
	})
}

func ToolsStatus(c *gin.Context) {
	sta := Sta_GetStatuses()
	c.HTML(http.StatusOK, "tools_status.html", gin.H{
		"appbase": appbase(c),
		"title":   app_conf.AppInfo().AppName + " - Statuses",
		"js":      "tools.js",
		"sta":     sta,
	})
}

func ToolsTypes(c *gin.Context) {
	typ := Typ_GetTypes()
	c.HTML(http.StatusOK, "tools_types.html", gin.H{
		"appbase": appbase(c),
		"title":   app_conf.AppInfo().AppName + " - Types",
		"js":      "tools.js",
		"typ":     typ,
	})
}
