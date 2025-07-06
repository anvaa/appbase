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

	tools_Api(r) // register tools API routes

	// SET app routes
	appGrp := r.Group("/app")
	{
		appGrp.Use(middleware.RequireAuth)

		appGrp.GET("/", app_ctrl.MainMenu)
		appGrp.GET("/start", app_ctrl.Start)
		
	}

	

	return r

}
