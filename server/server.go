package server

import (
	"app/app_api"
	"server/srv_conf"
	"server/middleware"

	"fmt"

	"github.com/gin-gonic/gin"

)

func InitWebServer() *gin.Engine {

	// GET gin mode from app.yaml
	if srv_conf.IsGinModDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Add security and CORS middleware
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			GinError(c)
		}
	})
	r.Use(middleware.CSRFProtection())

	// handle 500
	r.NoRoute(GinError)

	GinLoggerDatabase(r)

	setupRoutes(r)

	return r
}

func setupRoutes(r *gin.Engine) *gin.Engine {

	// Set up the app routes
	r = app_api.App_Api(r) // sets the routes for the app package

	return r
}

// send all error to /error page
func GinError(c *gin.Context) {
	errtxt := fmt.Sprintf("%v", c.Errors)
	fmt.Println("Error:", errtxt)
	c.HTML(500, "error.html", gin.H{
		"error": errtxt,
		"code":  c.Writer.Status(),
	})
}

// CORSMiddleware sets a strict CORS policy
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "https://yourdomain.com")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
