package middleware

import (
	"app/app_db"
	"server/srv_sec"
	"user/user_conf"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	userKey  = user_conf.UserKey
	UserUUID uint
	UserID   uint
)

// Verify middleware checks JWT and user authentication, used for both internal and external verification.
func Verify(c *gin.Context) {
	fmt.Println("Running Verify middleware")
	tokenString := getTokenFromRequest(c)
	if tokenString == "" {
		handleAuthFailure(c)
		return
	}
	fmt.Println("Token found:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(srv_sec.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		handleAuthFailure(c)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		handleAuthFailure(c)
		return
	}

	// Handle sub claim - it could be either string or number
	var sub string
	switch subClaim := claims["sub"].(type) {
	case string:
		sub = subClaim
	case float64:
		sub = fmt.Sprintf("%.0f", subClaim)
	case int:
		sub = fmt.Sprintf("%d", subClaim)
	default:
		log.Printf("Invalid subject claim type: %v (type: %T)", claims["sub"], claims["sub"])
		handleAuthFailure(c)
		return
	}
	if sub == "" {
		log.Printf("Empty subject claim: %v", claims["sub"])
		handleAuthFailure(c)
		return
	}

	user, err := app_db.User_GetByUUID(sub)
	if err != nil || user.ID == 0 || !user.IsAuth {
		handleAuthFailure(c)
		return
	}

	tokenVersion, ok := claims["token_version"].(float64)
	if !ok || int(tokenVersion) != user.TokenVersion {
		handleAuthFailure(c)
		return
	}

	c.Set(userKey, user)
	UserID = user.ID
	UserUUID = user.UUID

	// For /user/verify endpoint, respond with user data (SvelteKit compatibility)
	if c.FullPath() == "/verify" || c.FullPath() == "/user/verify" {
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Content-Type", "application/json")

		// Calculate expiration from JWT claims
		exp, _ := claims["exp"].(float64)
		expiresAt := time.Unix(int64(exp), 0).Format(time.RFC3339)

		c.JSON(http.StatusOK, gin.H{
			"username":   user.Username,
			"roles":      []string{user.AuthLevel.Name},
			"expires_at": expiresAt,
			"profile": gin.H{
				"uuid":       user.UUID,
				"auth_level": user.AuthLevel.Name,
				"orgs":       user.Org,
				"note":       user.Note,
			},
		})

		c.Abort()
		return
	}

	log.Printf("User %s (ID: %d) authenticated successfully", user.Username, user.ID)
	c.Next()
}

// getTokenFromRequest extracts JWT token from either Authorization header or cookie
func getTokenFromRequest(c *gin.Context) string {
	// First try Authorization header (Bearer token) - for SvelteKit/API clients
	auth := c.GetHeader("Authorization")
	if auth != "" {
		if len(auth) > 7 && auth[:7] == "Bearer " {
			return auth[7:] // Remove "Bearer " prefix
		}
	}

	// Then try cookie - for browser clients
	cookiename := user_conf.CookieName
	fmt.Println("Looking for cookie:", cookiename)
	tokenString, err := c.Cookie(cookiename)
	if err == nil && tokenString != "" {
		return tokenString
	}

	// Check for token in request body (JSON)
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err == nil {
		if token, ok := body["token"].(string); ok && token != "" {
			return token
		}
		if session, ok := body["session"].(string); ok && session != "" {
			return session
		}
	}

	return ""
}

// handleAuthFailure handles authentication failures differently for API vs browser requests
func handleAuthFailure(c *gin.Context) {
	// Check if this is an API request (has Authorization header or Accept: application/json)
	auth := c.GetHeader("Authorization")
	accept := c.GetHeader("Accept")

	if auth != "" || accept == "application/json" {
		// API client - return JSON error
		log.Println("Invalid or expired token")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		c.Abort()
		return
	}

	log.Println("Authentication required - redirecting to login")
	// Browser client - redirect to login
	redirectToLogin(c)
}

func redirectToLogin(c *gin.Context) {
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", false, true)
	c.Redirect(302, "/login")
	c.Abort()
}
