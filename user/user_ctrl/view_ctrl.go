package user_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"server/middleware"
	"server/srv_conf"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func Version(c *gin.Context) {
	c.String(http.StatusOK, "%s", app_conf.AppInfo())
}

func Info(c *gin.Context) {
	
	appinfo := app_conf.AppInfo()
	title := fmt.Sprintf("%s v%s", appinfo.AppName, appinfo.Version)

	c.HTML(http.StatusOK, "info.html", gin.H{
		"title":   title,
		"css":     "index.css",
		
		"appinfo": appinfo,
		"loctime": app_conf.GetLocalTime(),
		"apptime": app_conf.RunTime(),
	})
}

func View_Signup(c *gin.Context) {
	if c.Request.Method == "POST" {
		SignUp(c)
		return
	}
	count := "0"
	if c.Param("count") != "" {
		count = c.Param("count")
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{
		"title": "Signup",
		"css":   "index.css",
		"js":    "index.js",

		"count":      count,
		"logo_small": app_conf.AppLogos()[0],
	})
}

func View_Login(c *gin.Context) {
	if c.Request.Method == "POST" {
		Login(c)
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
		"css":   "index.css",
		"js":    "index.js",

		"logo_small": app_conf.AppLogos()[0],
		"logo_large": app_conf.AppLogos()[1],
	})
}

func View_NewUsers(c *gin.Context) {
	new_users, _ := app_db.Users_GetNew()

	c.HTML(http.StatusOK, "users_new.html", gin.H{
		"appbase":   middleware.AppBase(c),
		"title":     "New users",
		"js":        "users.js",
		"new_users": new_users,
	})
}

func View_ManageUsers(c *gin.Context) {

	auth_users, _ := app_db.Users_GetAuth()
	unauth_users, _ := app_db.Users_GetUnAuth()
	del_users, _ := app_db.Users_GetDeleted()

	c.HTML(http.StatusOK, "users.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   "Manage Users",
		"js":      "users.js",

		"auth_users":   auth_users,
		"unauth_users": unauth_users,
		"del_users":    del_users,
	})
}

func View_EditUser(c *gin.Context) {
	uuid := c.Param("uuid")

	edit_user, err := app_db.User_GetByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}

	auth_levels := app_db.GetAuthLevels()

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   "Edit User",
		"js":      "users.js",
		"css":     "tools.css",

		"edituid":     edit_user,
		"auth_levels": auth_levels,
	})
}

func View_MyAccount(c *gin.Context) {

	c.HTML(http.StatusOK, "myaccount.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   "My Account",
		"css":     "",
		"js":      "myaccount.js",
	})
}

func View_Database(c *gin.Context) {

	dbConfig := srv_conf.GetDbConfig()

	c.HTML(http.StatusOK, "db_config.html", gin.H{
		"appbase": middleware.AppBase(c),
		"title":   "Database Config",
		"css":     "",
		"js":      "db.js",

		"dbconf": dbConfig,
	})
}
