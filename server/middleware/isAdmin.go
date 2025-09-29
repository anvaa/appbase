package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/app_models"
)

// IsAdmin middleware checks if the authenticated user has admin privileges (Level >= 40)
// This middleware assumes that isAuth has already run and set the user context
func IsAdmin(c *gin.Context) {
	// Get user from context (should be set by isAuth middleware)
	userValue, exists := c.Get(userKey)
	if !exists {
		handleAuthorizationFailure(c, "Authentication required")
		return
	}

	user, ok := userValue.(app_models.Users)
	if !ok {
		handleAuthorizationFailure(c, "Invalid user context")
		return
	}

	// Check if user has admin privileges (Level >= 40)
	// Admin level: 40, Super level: 30, Manager level: 20, User level: 10, Guest level: 1
	if user.ID == 0 || user.AuthLevel.Level < 40 {
		handleAuthorizationFailure(c, "Admin privileges required")
		return
	}
	log.Println("IsAdmin passed for user:", user.Username, "with AuthLevel:", user.AuthLevel.Level)
	c.Next()
}

// IsSuper middleware checks if the authenticated user has super user privileges (Level >= 30)
// This allows both super users and admins to access the resource
func IsSuper(c *gin.Context) {
	// Get user from context (should be set by isAuth middleware)
	userValue, exists := c.Get(userKey)
	if !exists {
		handleAuthorizationFailure(c, "Authentication required")
		return
	}

	user, ok := userValue.(app_models.Users)
	if !ok {
		handleAuthorizationFailure(c, "Invalid user context")
		return
	}

	// Check if user has super user privileges (Level >= 30)
	// This allows both super users (Level 30) and admins (Level 40)
	if user.ID == 0 || user.AuthLevel.Level < 30 {
		handleAuthorizationFailure(c, "Super user privileges required")
		return
	}

	c.Next()
}

// IsManager middleware checks if the authenticated user has manager privileges (Level >= 20)
// This allows managers, super users, and admins to access the resource
func IsManager(c *gin.Context) {
	// Get user from context (should be set by isAuth middleware)
	userValue, exists := c.Get(userKey)
	if !exists {
		handleAuthorizationFailure(c, "Authentication required")
		return
	}

	user, ok := userValue.(app_models.Users)
	if !ok {
		handleAuthorizationFailure(c, "Invalid user context")
		return
	}

	// Check if user has manager privileges (Level >= 20)
	if user.ID == 0 || user.AuthLevel.Level < 20 {
		handleAuthorizationFailure(c, "Manager privileges required")
		return
	}

	c.Next()
}

// handleAuthorizationFailure handles authorization failures with context-aware responses
func handleAuthorizationFailure(c *gin.Context, message string) {
	// Check if this is an API request (has Authorization header or Accept: application/json)
	auth := c.GetHeader("Authorization")
	accept := c.GetHeader("Accept")

	if auth != "" || accept == "application/json" {
		// API client - return JSON error
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, gin.H{"message": message})
		c.Abort()
		return
	}

	// Browser client - return HTML error
	c.HTML(http.StatusForbidden, "error.html", gin.H{
		"error": "Forbidden: " + message,
	})
	c.Abort()
}
