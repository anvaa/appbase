package app_api

import (
	
	"srv/middleware"
	"srv/srv_conf"

	"app/app_ctrl"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routes for the application
func App_Api(r *gin.Engine) *gin.Engine {

	// Set up CSRF protection
	r.Use(middleware.CSRFProtection()) // set up store
	r.Use(middleware.CSRF()) // set up CSRF middleware
	r.Use(middleware.CSRFToken()) // add to context

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
	error_Api(r) // register error API routes
	user_Api(r) // register user API routes
	tools_Api(r) // register tools API routes

	// SET app routes
	appGrp := r.Group("/app")
	{
		appGrp.Use(middleware.IsAuth)

		appGrp.GET("/", app_ctrl.MainMenu)
		appGrp.GET("/start", app_ctrl.Start)
		
	}

	projGrp := r.Group("/proj")
	{
		projGrp.Use(middleware.IsAuth)

		projGrp.GET("/:id", app_ctrl.Proj_View) // view project page
		projGrp.GET("/edit/:id", app_ctrl.Proj_Edit)	 // edit project page
		
		projGrp.POST("/addupd", app_ctrl.Proj_AddUpd) // add or update project function
		projGrp.DELETE("/:id", app_ctrl.Proj_Delete) // delete project function
	
	}
	

	return r

}
