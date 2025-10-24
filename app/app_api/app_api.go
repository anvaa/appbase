package app_api

import (
	//"server/middleware"
	"server/srv_conf"

	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes for the application
func App_Api(r *gin.Engine) *gin.Engine {

	// SET app paths
	static_dir := srv_conf.StaticDir
	r.Static("/media", static_dir+"/media")
	r.Static("/css", static_dir+"/css")
	r.Static("/js", static_dir+"/js")
	r.Static("/assets", srv_conf.AssetsDir)

	r.LoadHTMLGlob(static_dir + "/html/*.html")

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("/media/favicon.ico")
	})
	r.GET("/robots.txt", func(c *gin.Context) {
		c.File("/media/robots.txt")
	})

	// Import API routes
	r = error_Api(r) // register error API routes
	r = user_Api(r)  // register user API routes
	r = tools_Api(r) // register tools API routes

	// SET app routes
	appGrp := r.Group("/app")
	{
		// appGrp.Use(middleware.Verify)
		// appGrp.Use(middleware.RequireLevel(10)) // user level

		appGrp.GET("/", app_ctrl.Start)

	}

	return r

}
