package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"server/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(c *gin.Context) {
	menu := app_db.Get_MenuTitles()
	c.HTML(http.StatusOK, "start.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   app_conf.AppInfo().AppName,
		"js":      "start.js",
		"menu0":   menu[0],
	})
}

func ToolsTitles(c *gin.Context) {
	menu := app_db.Get_MenuTitles()
	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   app_conf.AppInfo().AppName + " - Titles",
		"js":      "tools.js",
		"menu":    menu,
	})
}

func ToolsStatus(c *gin.Context) {
	sta := Sta_GetStatuses()
	c.HTML(http.StatusOK, "tools_status.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   app_conf.AppInfo().AppName + " - Statuses",
		"js":      "tools.js",
		"sta":     sta,
	})
}

func ToolsTypes(c *gin.Context) {
	typ := Typ_GetTypes()
	c.HTML(http.StatusOK, "tools_types.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   app_conf.AppInfo().AppName + " - Types",
		"js":      "tools.js",
		"typ":     typ,
	})
}
