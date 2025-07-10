package app_api
import (
	"github.com/gin-gonic/gin"
)

// On any GIN errer, redirect to the error page
func error_Api(r *gin.Engine) *gin.Engine {
	// Set up the error routes
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error.html", gin.H{
			"title": "Page Not Found",
			"css":   "index.css",
			"js":    "index.js",
			"error": "404 - Page Not Found",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.HTML(405, "error.html", gin.H{
			"title": "Method Not Allowed",
			"css":   "index.css",
			"js":    "index.js",
			"error": "405 - Method Not Allowed",
		})
	})

	return r
}