package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"net/http"
	"time"
	"fmt"

	"app/app_conf"
	"app/app_models"
	"srv/srv_sec"
)

func SetJWT(c *gin.Context, user *app_models.Users) (*gin.Context, error) {
	// Generate a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.UUID,
		"exp":  time.Now().Add(app_conf.SessionExpire()).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "appbase",
		"aud":  "appbase",
		"role": user.Role,
		"auth": user.IsAuth,
	})

	tokenString, err := token.SignedString([]byte(srv_sec.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate"})
		return nil, err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(app_conf.CookieName, tokenString, app_conf.CookieExpire, "/", "", true, true)

	setHeaders(c)
	setCookies(c)
	setSession(c, user)

	return c, nil
}

func setSession(c *gin.Context, user *app_models.Users) (*gin.Context, error) {

	expire := time.Now().Add(app_conf.SessionExpire())
	if expire.IsZero() {
		expire = time.Now().Add(24 * time.Hour) // Default to 24 hours if not set
	}
	if expire.Before(time.Now()) {
		return nil, gin.Error{
			Err:  fmt.Errorf("session expiration time is in the past"),
			Type: gin.ErrorTypePublic,
		}
	}

	local := c.GetHeader("Accept-Language")
	if local == "" {
		local = "en-US" // Default locale if not provided
	}

	// Create a new session with a unique ID and the user's information
	session := app_models.Session{
		SessionID: srv_sec.UUID_String(),
		UserID:    user.ID,
		UserUUID:  user.UUID,
		Locale:    local,
		Expire:    expire,
	}
	fmt.Println("Session created:", session)
	c.Set("session", session)
	c.Set("user", user)

	return c, nil
}

func setHeaders(c *gin.Context) {
	// Set default headers
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
	c.Next()
}

func setCookies(c *gin.Context) {
	// Set default cookies
	c.SetCookie("csrf_token", srv_sec.JwtSecret, app_conf.CookieExpire, "/", "", true, true)
	c.Next()
}