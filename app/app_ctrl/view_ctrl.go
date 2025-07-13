package app_ctrl

import (
	"app/app_conf"
	"app/app_db"

	"net/http"

	"github.com/gin-gonic/gin"
)

type api_base struct {
	Title   string
	User    interface{}
	Appinfo interface{}
	CSRF    string // CSRF token
}

var (
	appname = app_conf.AppName
	appinfo = app_conf.AppInfo()

	apibase = api_base{
		Title:   appname,
		User:    nil, // Will be set in the context
		Appinfo: appinfo,
		CSRF:    "", // Will be set in the context
	}
)

func SetApiBase(c *gin.Context) {
	apibase.Appinfo = appinfo
	apibase.User = c.Keys["user"]
	apibase.CSRF = c.MustGet("csrf").(string) // CSRF token
	c.Set("apibase", apibase)
}



func MainMenu(c *gin.Context) {
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"title":   app_conf.AppName,
		"user":    c.Keys["user"],
		"appinfo": appinfo,
	})
}

func Start(c *gin.Context) {

	projects := app_db.GetAllProjects()

	c.HTML(http.StatusOK, "start.html", gin.H{
		
		"title":   appname,
		"js":      "start.js",
		// "user":    c.Keys["user"],
		// "appinfo": appinfo,
		// "csrf":    c.MustGet("csrf").(string), // CSRF token

		"projects": projects,
	})
}

func ToolsTitles(c *gin.Context) {

	menu := app_db.Get_MenuTitles()

	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"title":   app_conf.AppName + " - Titles",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"menu": menu,

	})
}

func ToolsStatus(c *gin.Context) {

	sta := Sta_GetStatuses()

	c.HTML(http.StatusOK, "tools_status.html", gin.H{
		"title":   app_conf.AppName + " - Statuses",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"sta": sta,

	})
}

func ToolsTypes(c *gin.Context) {

	typ := Typ_GetTypes()

	c.HTML(http.StatusOK, "tools_types.html", gin.H{
		"title":   app_conf.AppName + " - Types",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"typ": typ,
	})
}
