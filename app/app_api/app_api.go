package app_api

import (
	"srv/srv_conf"
	"srv/web/middleware"

	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes for the application
func App_Api(r *gin.Engine) *gin.Engine {

	// SET app paths
	static_dir := srv_conf.StaticDir

	r.Static("/css", static_dir+"/css")
	r.Static("/js", static_dir+"/js")
	// r.Static("/media", static_dir+"/media")

	r.Static("/assets", srv_conf.AssetsDir)
	

	r.LoadHTMLGlob(static_dir + "/html/*.html")

	// SET app routes
	invGrp := r.Group("/app")
	{
		invGrp.Use(middleware.RequireAuth)

		invGrp.GET("/", app_ctrl.MainMenu)
		invGrp.GET("/start", app_ctrl.Start)
		

	}

	manGrp := r.Group("/subitem")
	{
		manGrp.Use(middleware.RequireAuth)

		manGrp.POST("/addupd", app_ctrl.Sub_AddUpd)
		manGrp.POST("/delete", app_ctrl.Sub_Delete)

	}

	staGrp := r.Group("/status")
	{
		staGrp.Use(middleware.RequireAuth)

		staGrp.POST("/addupd", app_ctrl.Sta_AddUpd)
		staGrp.POST("/delete", app_ctrl.Sta_Delete)

	}

	titelsGrp := r.Group("/title")
	{
		titelsGrp.Use(middleware.RequireAuth)

		titelsGrp.POST("/upd", app_ctrl.Mnu_UpdTitels)

	}

	toolsGrp := r.Group("/tools")
	{
		toolsGrp.Use(middleware.RequireAuth)

		toolsGrp.GET("/titles", app_ctrl.ToolsTitles)
		toolsGrp.GET("/conf", app_ctrl.ToolsConf) // tools configuration page

		toolsGrp.POST("/printconf", app_ctrl.PrintConf) // save print conf
		toolsGrp.POST("/appconf", app_ctrl.AppConf) //  application conifg/settings

	}

	return r

}
