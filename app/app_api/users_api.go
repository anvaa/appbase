package app_api

import (

	"github.com/gin-gonic/gin"

	"server/middleware"
	"user/user_ctrl"
	"user"
)

func user_Api(r *gin.Engine) *gin.Engine {

	// SET default paths
	r.GET("/", user.Root)
	r.GET("/info", user.Info)
	r.GET("/version", user.Version)

	// Set up the user routes
	r.POST("/signup", user.View_Signup)
	r.GET("/signup/:count", user.View_Signup)
	r.POST("/login", user.View_Login)
	r.GET("/login", user.View_Login)

	r.GET("/logout", middleware.Logout)

	userRoutes := r.Group("/user")
	{
		userRoutes.Use(middleware.IsAuth)

		userRoutes.POST("/psw", user_ctrl.User_SetNewPassword)

		userRoutes.Use(middleware.IsAdmin)

		userRoutes.GET("/", user_ctrl.GetAllUsers)
		userRoutes.GET("/:id", user_ctrl.GetUser)

		userRoutes.POST("/delete", user_ctrl.User_DeleteUser)
		userRoutes.POST("/auth", user_ctrl.User_UpdateAuth)
		userRoutes.POST("/role", user_ctrl.User_UpdateRole)
		userRoutes.POST("/org", user_ctrl.User_UpdateOrg)
	}

	viewRoutes := r.Group("/v")
	{
		viewRoutes.Use(middleware.IsAuth)

		// is users
		viewRoutes.GET("/myaccount", user.View_MyAccount)

		viewRoutes.Use(middleware.IsAdmin)

		// is admin
		viewRoutes.GET("/newusers", user.View_NewUsers)
		viewRoutes.GET("/users", user.View_ManageUsers)
		viewRoutes.GET("/user/:uid", user.View_EditUser)

	}

	return r
}
