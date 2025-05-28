package user_api

import (

	"log"

	"github.com/gin-gonic/gin"

	"srv/web/middleware"
	"srv/srv_conf"

	"usr"

)

func User_Api(r *gin.Engine) *gin.Engine {

	// SET app paths
	static_dir := srv_conf.StaticDir
	log.Println("Static App folder:", static_dir)
	r.Static("/media", static_dir+"/media")

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("/media/favicon.ico")
	})
	r.GET("/robots.txt", func(c *gin.Context) {
		c.File("/media/robots.txt")
	})

	// SET user paths
	r.GET("/", users.Root)
	r.GET("/info", users.Info)
	r.GET("/version", users.Version)

	r.POST("/signup", users.View_Signup)
	r.GET("/signup/:count", users.View_Signup)
	r.POST("/login", users.View_Login)
	r.GET("/login", users.View_Login)
	
	r.GET("/logout", middleware.Logout)

	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.RequireAuth)

		userRoutes.POST("/psw", users.User_SetNewPassword)

		
		userRoutes.Use(middleware.IsAdmin)

		userRoutes.GET("/", users.GetAllUsers)
		userRoutes.GET("/:id", users.GetUser)

		userRoutes.POST("/delete", users.User_DeleteUser)
		userRoutes.POST("/auth", users.User_UpdateAuth)
		userRoutes.POST("/role", users.User_UpdateRole)
		userRoutes.POST("/org", users.User_UpdOrg)
	}

	viewRoutes := r.Group("/v")
	{
		viewRoutes.Use(middleware.RequireAuth)

		// is users
		viewRoutes.GET("/myaccount", users.View_MyAccount)

		viewRoutes.Use(middleware.IsAdmin)

		// is admin
		viewRoutes.GET("/newusers", users.View_NewUsers)
		viewRoutes.GET("/users", users.View_ManageUsers)
		viewRoutes.GET("/user/:edituid", users.View_EditUser)

	}

	return r
}
