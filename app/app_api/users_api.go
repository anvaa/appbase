package app_api

import (

	"github.com/gin-gonic/gin"

	"server/middleware"
	"user/user_ctrl"

)

func user_Api(r *gin.Engine) *gin.Engine {
	
	// root routes
	rootGrp := r.Group("/")
	{
		rootGrp.GET("/", user_ctrl.Root)

		rootGrp.GET("/info", user_ctrl.Info)
		rootGrp.GET("/health", user_ctrl.Health)

		rootGrp.POST("/signup", user_ctrl.View_Signup)
		rootGrp.GET("/signup/:count", user_ctrl.View_Signup)
		
		rootGrp.POST("/login", user_ctrl.Login)
		rootGrp.GET("/login", user_ctrl.View_Login)

		rootGrp.Use(middleware.IsAuth)
		rootGrp.POST("/verify", middleware.Verify)
		rootGrp.GET("/logout", middleware.Logout)
	}

	// User routes
	userGrp := r.Group("/user")
	{
		userGrp.Use(middleware.IsAuth)

		userGrp.POST("/psw", user_ctrl.User_SetNewPassword)

		userGrp.Use(middleware.IsAdmin)

		userGrp.GET("/", user_ctrl.GetAllUsers)
		userGrp.GET("/:id", user_ctrl.GetUser)

		userGrp.POST("/delete", user_ctrl.User_DeleteUser)
		userGrp.POST("/auth", user_ctrl.User_UpdateAuth)
		userGrp.POST("/authlevel", user_ctrl.User_UpdAuthLevel)
		userGrp.POST("/org", user_ctrl.User_UpdateOrg)
	}

	// View routes
	viewGrp := r.Group("/v")
	{
		viewGrp.Use(middleware.IsAuth)

		// is users
		viewGrp.GET("/myaccount", user_ctrl.View_MyAccount)

		viewGrp.Use(middleware.IsAdmin)

		// is admin
		viewGrp.GET("/newusers", user_ctrl.View_NewUsers)
		viewGrp.GET("/users", user_ctrl.View_ManageUsers)
		viewGrp.GET("/user/:uuid", user_ctrl.View_EditUser)

		viewGrp.GET("/orgs", user_ctrl.Org_View)
		viewGrp.GET("/org/new", user_ctrl.Org_New)

		viewGrp.GET("/org/:uuid", user_ctrl.Org_Edit)
		viewGrp.POST("/org/addupd", user_ctrl.Org_AddUpd)
		viewGrp.DELETE("/org/:uuid", user_ctrl.Org_Delete)

		viewGrp.GET("/org/members/:uuid", user_ctrl.Org_Members)
		viewGrp.POST("/org/members/add", user_ctrl.Org_AddMember)
		viewGrp.POST("/org/members/rem", user_ctrl.Org_RemoveMember)

		// is database
		viewGrp.GET("/database", user_ctrl.View_Database)
		viewGrp.POST("/dbconf", user_ctrl.DB_SaveDbConf)

	}

	return r
}
