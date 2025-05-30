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
	stats, c_stats := Stat_GetStatusSubItems()

	c.HTML(http.StatusOK, "tools_titles.html", gin.H{
		"title":   app_conf.AppName + " - Titles",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"menu":  menu,
		"menu0": menu[0],
		"menu1": menu[1],
		"menu2": menu[2],
		"menu3": menu[3],
		"menu4": menu[4],
		"menu5": menu[5],
		"menu6": menu[6],
		"menu7": menu[7],
		"menu8": menu[8],

		"stats":   stats,
		"c_stats": c_stats,
	})
}

func ToolsConf(c *gin.Context) {

	printConf := app_conf.PrintConf()
	appConf := app_conf.AppConf()

	c.HTML(http.StatusOK, "tools_conf.html", gin.H{
		"title":   app_conf.AppName + " - Config",
		"js":      "tools.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"printConf": printConf,
		"appConf":   appConf,
	})
}
