package middleware

import (
	"github.com/gin-gonic/gin"

	"app/app_models"
)

func IsAdmin(c *gin.Context) {

	user := c.MustGet("user").(app_models.Users)
	if user.Role != "admin" {
		OnErr(c)
		return
	}
	c.Next()
}
