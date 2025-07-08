package middleware

import (

	"github.com/gin-gonic/gin"

	"app/app_models"

)

func IsAuth(c *gin.Context) {

	user := c.MustGet("user").(app_models.Users)
	if user.ID == 0 || !user.IsAuth {
		c.HTML(401, "error.html", gin.H{
			"error": "Unauthorized access. Please log in.",
		})
		c.Abort() // Stop further processing
		return
	}
	// User is authenticated, proceed with the request
	c.Next()
}
