package users

import (
	"app/app_conf"

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
		"url":     c.Request.Referer(),
		"info":    appinfo,
		"company": appinfo.AppName,
		"loctime": app_conf.GetLocalTime(),
		"apptime": app_conf.RunTime(),
		"backbtn": c.Request.Referer(),
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
		"count": count,
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
	})
}

func View_NewUsers(c *gin.Context) {
	newusers, cnewusers, _ := Users_GetNew()

	c.HTML(http.StatusOK, "users_new.html", gin.H{
		"title":    "New users",
		"user":     c.Keys["user"],
		"js":       "users.js",
		"newusers": newusers,
		"countnew": cnewusers,
	})
}

func View_ManageUsers(c *gin.Context) {
	fmt.Println("View_ManageUsers called")
	auth_users, c_auth, _ := Users_GetAuth()
	fmt.Println("auth_users count:", c_auth)
	unauth_users, c_unauth, _ := Users_GetUnAuth()
	fmt.Println("unauth_users count:", c_unauth)
	del_users, c_del, _ := Users_GetDeleted()
	fmt.Println("del_users count:", c_del)

	c.HTML(http.StatusOK, "users.html", gin.H{
		"title": "Manage Users",
		"user":  c.Keys["user"],
		"js":    "users.js",

		"authusers":   auth_users,
		"countauth":   c_auth,
		"unauthusers": unauth_users,
		"countunauth": c_unauth,
		"delusers":    del_users,
		"countdel":    c_del,
	})
}

func View_EditUser(c *gin.Context) {
	edit_id := c.Param("uid")

	edit_user, err := User_GetById(edit_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}

	c.HTML(http.StatusOK, "user_edit.html", gin.H{
		"title": "Edit User",
		"js":    "users.js",
		"user":  c.Keys["user"],

		"edituid": edit_user,
	})
}

func View_MyAccount(c *gin.Context) {

	c.HTML(http.StatusOK, "myaccount.html", gin.H{
		"title": "My Account",
		"css":   "",
		"js":    "myaccount.js",
		"user":  c.Keys["user"],
	})
}
