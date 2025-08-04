package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"srv/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appname = app_conf.AppName
	appinfo = app_conf.AppInfo()
	CurUserID int
	CurUserUUID int
)

type api_base struct {
	User    any
	Appinfo any
}

func setApiBase(c *gin.Context) any {
	CurUserUUID = middleware.UserUUID
	CurUserID = middleware.UserID
	apibase := api_base{
		User:    c.Keys[app_conf.UserKey],
		Appinfo: appinfo,
	}
	return apibase
}

func MainMenu(c *gin.Context) {
	middleware.ValidateCSRFToken(c)
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"apibase": setApiBase(c),
		"title":   app_conf.AppName,
	})
}

func Start(c *gin.Context) {
	middleware.ValidateCSRFToken(c)
	projects := app_db.GetAllProjects()

	c.HTML(http.StatusOK, "start.html", gin.H{
		"apibase": setApiBase(c),

		"title": appname,
		"js":    "start.js",
		"css":   "tree.css",

		"projects": projects,
	})
}

func ToolsTitles(c *gin.Context) {
	middleware.ValidateCSRFToken(c)
	menu := app_db.Get_MenuTitles()

	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"apibase": setApiBase(c),
		"title":   app_conf.AppName + " - Titles",
		"js":      "tools.js",

		"menu": menu,
	})
}

func ToolsStatus(c *gin.Context) {
	middleware.ValidateCSRFToken(c)
	sta := Sta_GetStatuses()

	c.HTML(http.StatusOK, "tools_status.html", gin.H{
		"apibase": setApiBase(c),
		"title":   app_conf.AppName + " - Statuses",
		"js":      "tools.js",

		"sta": sta,
	})
}

func ToolsTypes(c *gin.Context) {
	middleware.ValidateCSRFToken(c)
	typ := Typ_GetTypes()

	c.HTML(http.StatusOK, "tools_types.html", gin.H{
		"apibase": setApiBase(c),
		"title":   app_conf.AppName + " - Types",
		"js":      "tools.js",

		"typ": typ,
	})
}
