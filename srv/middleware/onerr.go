package middleware

import (
	"log"
	"app/app_conf"

	"github.com/gin-gonic/gin"

	"net/http"
)

func OnErr(c *gin.Context) {

	// error handling here
	log.Println("Error: ", c.Errors.String())

	Logout(c)
}

func Logout(c *gin.Context) {
	// delete cookie from browser, redirect to login page
	c.SetCookie(app_conf.CookieName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/")

	
}
