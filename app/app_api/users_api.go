package app_api

import (
	"github.com/gin-gonic/gin"

	"server/middleware"
	"user/user_ctrl"
)

func UserAPI(r *gin.Engine) {
	registerRootRoutes(r)
	registerUserRoutes(r)
	registerViewRoutes(r)
}

func registerRootRoutes(r *gin.Engine) {
	rootGrp := r.Group("/")
	rootGrp.GET("/", user_ctrl.Root)
	rootGrp.GET("/info", user_ctrl.Info)
	rootGrp.GET("/health", user_ctrl.Health)
	rootGrp.POST("/signup", user_ctrl.View_Signup)
	rootGrp.GET("/signup/:count", user_ctrl.View_Signup)
	rootGrp.POST("/login", user_ctrl.Login)
	rootGrp.GET("/login", user_ctrl.View_Login)
	rootGrp.GET("/logout", user_ctrl.Logout)
	rootGrp.POST("/verify", middleware.Verify)
	rootGrp.GET("/verify", middleware.Verify)
}

func registerUserRoutes(r *gin.Engine) {
	userGrp := r.Group("/user", middleware.Verify)
	userGrp.POST("/psw", user_ctrl.User_SetNewPassword)

	adminGrp := userGrp.Group("/", middleware.RequireLevel(30))
	adminGrp.GET("/", user_ctrl.GetAllUsers)
	adminGrp.GET("/:id", user_ctrl.GetUser)
	adminGrp.POST("/delete", user_ctrl.User_DeleteUser)
	adminGrp.POST("/auth", user_ctrl.User_UpdateAuth)
	adminGrp.POST("/authlevel", user_ctrl.User_UpdAuthLevel)
	adminGrp.POST("/org", user_ctrl.User_UpdateOrg)
}

func registerViewRoutes(r *gin.Engine) {
	viewGrp := r.Group("/v", middleware.Verify)
	viewGrp.GET("/myaccount", user_ctrl.View_MyAccount)

	superGrp := viewGrp.Group("/", middleware.RequireLevel(30))
	superGrp.GET("/newusers", user_ctrl.View_NewUsers)
	superGrp.GET("/users", user_ctrl.View_ManageUsers)
	superGrp.GET("/user/:uuid", user_ctrl.View_EditUser)
	superGrp.GET("/orgs", user_ctrl.Org_View)
	superGrp.GET("/org/new", user_ctrl.Org_New)
	superGrp.GET("/org/:uuid", user_ctrl.Org_Edit)
	superGrp.POST("/org/addupd", user_ctrl.Org_AddUpd)
	superGrp.DELETE("/org/:uuid", user_ctrl.Org_Delete)
	superGrp.GET("/org/members/:uuid", user_ctrl.Org_Members)
	superGrp.POST("/org/members/add", user_ctrl.Org_AddMember)
	superGrp.POST("/org/members/rem", user_ctrl.Org_RemoveMember)
	superGrp.GET("/database", user_ctrl.View_Database)
	superGrp.POST("/dbconf", user_ctrl.DB_SaveDbConf)
}
