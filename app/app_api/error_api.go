package app_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderErrorPage(c *gin.Context, status int, title, errorMsg string) {
	c.HTML(status, "error.html", gin.H{
		"title": title,
		"css":   "index.css",
		"js":    "index.js",
		"error": errorMsg,
	})
}

// ErrorApi sets up error handlers for 404 and 405 errors.
func ErrorApi(r *gin.Engine) *gin.Engine {
	r.NoRoute(func(c *gin.Context) {
		renderErrorPage(c, http.StatusNotFound, "Page Not Found", "404 - Page Not Found")
	})

	r.NoMethod(func(c *gin.Context) {
		renderErrorPage(c, http.StatusMethodNotAllowed, "Method Not Allowed", "405 - Method Not Allowed")
	})

	return r
}