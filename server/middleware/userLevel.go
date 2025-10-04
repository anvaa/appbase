package middleware

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"

    "app/app_models"
)

// RequireLevel returns a middleware function that checks if user has minimum required access level
func RequireLevel(minLevel int) gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Printf("Running RequireLevel middleware (minimum level: %d)\n", minLevel)
        
        userValue, exists := c.Get(userKey)
        if !exists {
            fmt.Println("DEBUG: No user found in context")
            handleAuthorizationFailure(c, "Authentication required")
            return
        }

        user, ok := userValue.(app_models.Users)
        if !ok {
            fmt.Printf("DEBUG: Type assertion failed. Type: %T\n", userValue)
            handleAuthorizationFailure(c, "Invalid user context")
            return
        }

        // Debug output
        fmt.Printf("DEBUG: User - ID: %d, Username: %s, AuthLevelID: %d\n", 
            user.ID, user.Username, user.AuthLevelID)
        
        if user.AuthLevel.ID != 0 {
            fmt.Printf("DEBUG: AuthLevel - ID: %d, Level: %d, Name: %s\n", 
                user.AuthLevel.ID, user.AuthLevel.Level, user.AuthLevel.Name)
        } else {
            fmt.Println("DEBUG: AuthLevel not loaded")
        }

        // Check if user exists and has minimum required level
        if user.ID == 0 {
            handleAuthorizationFailure(c, "Invalid user")
            return
        }

        // Use AuthLevel.Level for proper hierarchy checking
        if user.AuthLevel.Level < minLevel {
            fmt.Printf("DEBUG: Access denied - User level: %d, Required: %d\n", 
                user.AuthLevel.Level, minLevel)
            handleAuthorizationFailure(c, fmt.Sprintf("Insufficient privileges (required level: %d)", minLevel))
            return
        }

        // Log successful access
        log.Printf("Access granted to user: %s (Level: %d/%d)", 
            user.Username, user.AuthLevel.Level, minLevel)
        
        c.Next()
    }
}

// RequireRole returns a middleware function that checks for specific AuthLevelID
func RequireRole(roleID int) gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Printf("Running RequireRole middleware (role ID: %d)\n", roleID)
        
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

        // Check if user has the exact required role
        if user.ID == 0 || user.AuthLevelID != roleID {
            fmt.Printf("DEBUG: Role access denied - User role: %d, Required: %d\n", 
                user.AuthLevelID, roleID)
            handleAuthorizationFailure(c, fmt.Sprintf("Role required: %d", roleID))
            return
        }

        log.Printf("Role access granted to user: %s (Role: %d)", user.Username, user.AuthLevelID)
        c.Next()
    }
}

// Convenience middleware functions for common access levels
func RequireAdmin() gin.HandlerFunc {
    return RequireLevel(40) // Admin level
}

func RequireSuper() gin.HandlerFunc {
    return RequireLevel(30) // Super level or higher
}

func RequireManager() gin.HandlerFunc {
    return RequireLevel(20) // Manager level or higher
}

func RequireUser() gin.HandlerFunc {
    return RequireLevel(10) // User level or higher
}

// IsAdmin - specific role check (backwards compatibility)
func IsAdmin(c *gin.Context) {
    RequireRole(1)(c) // AuthLevelID = 1 for admin
}

// IsSuper - specific role check (backwards compatibility)  
func IsSuper(c *gin.Context) {
    RequireRole(2)(c) // AuthLevelID = 2 for super
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