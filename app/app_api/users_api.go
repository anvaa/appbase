package app_api

import (

	"github.com/gin-gonic/gin"

	"srv/middleware"

	"usr"
)

func user_Api(r *gin.Engine) *gin.Engine {

	// SET default paths
	r.GET("/", users.Root)
	r.GET("/info", users.Info)
	r.GET("/version", users.Version)

	// Set up the user routes
	r.POST("/signup", users.View_Signup)
	r.GET("/signup/:count", users.View_Signup)
	r.POST("/login", users.View_Login)
	r.GET("/login", users.View_Login)
	
	r.GET("/logout", middleware.Logout)

	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.IsAuth)

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
		viewRoutes.Use(middleware.IsAuth)

		// is users
		viewRoutes.GET("/myaccount", users.View_MyAccount)

		viewRoutes.Use(middleware.IsAdmin)

		// is admin
		viewRoutes.GET("/newusers", users.View_NewUsers)
		viewRoutes.GET("/users", users.View_ManageUsers)
		viewRoutes.GET("/user/:uid", users.View_EditUser)

	}

	return r
}
