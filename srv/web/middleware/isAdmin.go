package middleware

import (
	"github.com/gin-gonic/gin"

	"app/app_models"
)

func IsAdmin(c *gin.Context) {

	user := c.MustGet("user").(*app_models.Users)
	if user == nil || user.Role != "admin" {
		c.HTML(403, "error.html", gin.H{
			"error": "Forbidden: You do not have permission to access this resource.",
		})
		c.Abort() // Stop further processing
		return
	}
	// User is an admin, proceed with the request
	c.Next()
}

func IsSuper(c *gin.Context) {

	user := c.MustGet("user").(*app_models.Users)
	if user == nil || user.Role != "superuser" {
		c.HTML(403, "error.html", gin.H{
			"error": "Forbidden: You do not have permission to access this resource.",
		})
		c.Abort() // Stop further processing
		return
	}
	// User is an admin, proceed with the request
	c.Next()
}
