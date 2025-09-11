package app_api

import (

	"github.com/gin-gonic/gin"

	"server/middleware"
	"user/user_ctrl"

)

func user_Api(r *gin.Engine) *gin.Engine {

	// SET default paths
	r.GET("/", user_ctrl.Root)
	r.GET("/info", user_ctrl.Info)
	r.GET("/version", user_ctrl.Version)

	// Set up the user routes
	r.POST("/signup", user_ctrl.View_Signup)
	r.GET("/signup/:count", user_ctrl.View_Signup)
	r.POST("/login", user_ctrl.Login)
	r.GET("/login", user_ctrl.View_Login)

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
		userRoutes.POST("/authlevel", user_ctrl.User_UpdAuthLevel)
		userRoutes.POST("/org", user_ctrl.User_UpdateOrg)
	}

	viewRoutes := r.Group("/v")
	{
		viewRoutes.Use(middleware.IsAuth)

		// is users
		viewRoutes.GET("/myaccount", user_ctrl.View_MyAccount)

		viewRoutes.Use(middleware.IsAdmin)

		// is admin
		viewRoutes.GET("/newusers", user_ctrl.View_NewUsers)
		viewRoutes.GET("/users", user_ctrl.View_ManageUsers)
		viewRoutes.GET("/user/:uuid", user_ctrl.View_EditUser)

		viewRoutes.GET("/orgs", user_ctrl.Org_View)
		viewRoutes.GET("/org/new", user_ctrl.Org_New)

		viewRoutes.GET("/org/:uuid", user_ctrl.Org_Edit)
		viewRoutes.POST("/org/addupd", user_ctrl.Org_AddUpd)
		viewRoutes.DELETE("/org/:uuid", user_ctrl.Org_Delete)

		viewRoutes.GET("/org/members/:uuid", user_ctrl.Org_Members)
		viewRoutes.POST("/org/members/add", user_ctrl.Org_AddMember)
		viewRoutes.POST("/org/members/rem", user_ctrl.Org_RemoveMember)

		// is database
		viewRoutes.GET("/database", user_ctrl.View_Database)
		viewRoutes.POST("/dbconf", user_ctrl.DB_SaveDbConf)

	}

	return r
}
