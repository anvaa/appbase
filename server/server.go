package server

import (
	"app/app_api"
	"server/global"
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

func SaveHostInfo() error {
	// Save server info to srv_conf
	_ips, err := global.GetIPv4Addresses()
	if err != nil {
		return err
	}

	hostname := global.GetHostName()

	srv_conf.SetVal("host_ip", _ips)
	srv_conf.SetVal("host_name", hostname)

	return nil
}