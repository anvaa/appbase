package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	appinfo      = app_conf.AppInfo()
	
)



func Start(c *gin.Context) {

	c.HTML(http.StatusOK, "start.html", gin.H{
		"title": app_conf.AppName,
		"css": "start.css",
		"js":  "start.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		
	})
}

func Tools(c *gin.Context) {

	menu := app_db.Get_MenuTitles()
	printConf := app_conf.PrintConf()
	appConf := app_conf.AppConf()

	stats, c_stats := Stat_GetStatusSubItems()

	c.HTML(http.StatusOK, "tools.html", gin.H{
		"title": app_conf.AppName + " - Tools",
		"css": "tools.css",
		"js":  "start.js",
		"user":    c.Keys["user"],
		"appinfo": appinfo,

		"menu":    menu,
		"menu0":   menu[0],
		"menu1":   menu[1],
		"menu2":   menu[2],
		"menu3":   menu[3],
		"menu4":   menu[4],
		"menu5":   menu[5],
		"menu6":   menu[6],
		"menu7":   menu[7],
		"menu8":   menu[8],

		"printConf": printConf,
		"appConf":   appConf,

		"stats":   stats,
		"c_stats": c_stats,
	})
}