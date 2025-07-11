package app_ctrl

import (
	"app/app_conf"
	"app/app_db"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appinfo = app_conf.AppInfo()
)

func MainMenu(c *gin.Context) {
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"title":   app_conf.AppName,
		"user":    c.Keys["user"],
		"appinfo": appinfo,
	})
}

func Start(c *gin.Context) {

	menu := app_db.Get_MenuTitles()

	c.HTML(http.StatusOK, "start.html", gin.H{
		"title":   app_conf.AppName,
		"js":      "start.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"menu0": menu[0],
	})
}

func ToolsTitles(c *gin.Context) {

	menu := app_db.Get_MenuTitles()
	//sta := Sta_GetStatuses()
	//typ := Typ_GetTypes()

	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"title":   app_conf.AppName + " - Titles",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"menu":  menu,
		//"sta":  sta,
		//"typ":  typ,

	})
}

func ToolsStatus(c *gin.Context) {
	
	//menu := app_db.Get_MenuTitles()
	sta := Sta_GetStatuses()
	//typ := Typ_GetTypes()

	c.HTML(http.StatusOK, "tools_status.html", gin.H{
		"title":   app_conf.AppName + " - Statuses",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		//"menu":  menu,
		"sta":  sta,
		//"typ":  typ,
	})
}

func ToolsTypes(c *gin.Context) {

	//menu := app_db.Get_MenuTitles()
	//sta := Sta_GetStatuses()
	typ := Typ_GetTypes()

	c.HTML(http.StatusOK, "tools_types.html", gin.H{
		"title":   app_conf.AppName + " - Types",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		//"menu":  menu,
		//"sta":  sta,
		"typ":  typ,
	})
}
