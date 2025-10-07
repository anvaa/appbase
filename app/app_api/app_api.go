package app_api

import (
	"server/middleware"
	"server/srv_conf"
	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

// App_Api sets up the routes for the application
func App_Api(r *gin.Engine) *gin.Engine {
	staticDir := srv_conf.StaticDir

	// Static file routes
	r.Static("/media", staticDir+"/media")
	r.Static("/css", staticDir+"/css")
	r.Static("/js", staticDir+"/js")
	r.Static("/assets", srv_conf.AssetsDir)

	r.LoadHTMLGlob(staticDir + "/html/*.html")

	// Register API routes
	ErrorApi(r)
	UserAPI(r)
	ToolsApi(r)

	// App routes group with middleware
	appGrp := r.Group("/app", middleware.Verify)
	{
		appGrp.GET("/", app_ctrl.Start)
	}

	return r
}
