package middleware

import (
	"app/app_db"
	"app/app_models"
	"server/srv_conf"
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
	fmt.Println("Starting verification middleware")
	tokenString := extractToken(c)
	if tokenString == "" {
		authFailure(c)
		return
	}
	fmt.Println("Token found, length:", len(tokenString))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("ERROR: Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(srv_sec.JWTSecret()), nil
	})
	if err != nil {
		fmt.Printf("ERROR: JWT parsing failed: %v\n", err)
		authFailure(c)
		return
	}
	if !token.Valid {
		fmt.Printf("ERROR: Token is not valid\n")
		authFailure(c)
		return
	}
	fmt.Println("Token parsed successfully")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		authFailure(c)
		return
	}
	fmt.Println("Claims:", claims)

	sub, valid := parseSubClaim(claims["sub"])
	if !valid {
		authFailure(c)
		return
	}

	user, err := app_db.User_GetByUUID(sub)
	if err != nil || user.ID == 0 || !user.IsAuth {
		authFailure(c)
		return
	}

	if !validateTokenVersion(claims, user.TokenVersion) {
		authFailure(c)
		return
	}

	c.Set(userKey, user)
	UserID = user.ID
	UserUUID = user.UUID

	if isVerifyEndpoint(c.FullPath()) {
		respondWithUser(c, claims, user)
		return
	}

	fmt.Printf("User %s (ID: %d) authenticated successfully", user.Username, user.ID)
	c.Next()
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}

	tokenString, err := c.Cookie(user_conf.CookieName)
	if err == nil && tokenString != "" {
		return tokenString
	}

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

func parseSubClaim(subClaim interface{}) (string, bool) {
	switch v := subClaim.(type) {
	case string:
		if v != "" {
			return v, true
		}
	case float64:
		return fmt.Sprintf("%.0f", v), true
	case int:
		return fmt.Sprintf("%d", v), true
	}
	log.Printf("Invalid or empty subject claim: %v (type: %T)", subClaim, subClaim)
	return "", false
}

func validateTokenVersion(claims jwt.MapClaims, expected int) bool {
	tokenVersion, ok := claims["token_version"].(float64)
	return ok && int(tokenVersion) == expected
}

func isVerifyEndpoint(path string) bool {
	return path == "/verify" || path == "/user/verify"
}

func respondWithUser(c *gin.Context, claims jwt.MapClaims, user app_models.Users) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Header("Content-Type", "application/json")

	exp, _ := claims["exp"].(float64)
	expiresAt := time.Unix(int64(exp), 0).Format(time.RFC3339)

	c.JSON(http.StatusOK, gin.H{
		"username":   user.Username,
		"email":      user.Email,
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
}

func authFailure(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	accept := c.GetHeader("Accept")

	if auth != "" || accept == "application/json" {
		log.Println("Invalid or expired token")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		c.Abort()
		return
	}

	fmt.Println("Authentication required - redirecting to login")
	redirectToLogin(c)
}

func redirectToLogin(c *gin.Context) {
	// Clear cookie with the same secure setting used when it was set
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", srv_conf.UseTLS(), true)
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}
