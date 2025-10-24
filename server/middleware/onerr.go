package middleware

import (
	"log"
	"server/srv_conf"
	"user/user_conf"

	"github.com/gin-gonic/gin"

	"net/http"
)

func OnErr(c *gin.Context) {

	// error handling here
	log.Println("Error: ", c.Errors.String())

	Logout(c)
}

func Logout(c *gin.Context) {
	// Clear cookie for browser clients with the same secure setting used when it was set
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", srv_conf.UseTLS(), true)

	// Check if this is an API request
	accept := c.GetHeader("Accept")
	auth := c.GetHeader("Authorization")

	if accept == "application/json" || auth != "" {
		// API client - return JSON response
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
		return
	}

	// Browser client - redirect
	c.Redirect(http.StatusMovedPermanently, "/")
}
