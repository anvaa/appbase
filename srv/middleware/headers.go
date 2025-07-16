package middleware

import (
	
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"app/app_conf"
	"app/app_models"
	"srv/srv_sec"
)

// SetJWT generates a JWT, sets it as a cookie, and attaches session/user info to the context.
func SetJWT(c *gin.Context, user *app_models.Users) error {
	tokenString, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(app_conf.CookieName, tokenString, app_conf.CookieExpire, "/", "", true, true)

	setSecurityHeaders(c)
	setCSRFCookie(c)
	// if err := setSession(c, user); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set session"})
	// 	return err
	// }

	return nil
}

// generateJWT creates a signed JWT string for the given user.
func generateJWT(user *app_models.Users) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.UUID,
		"exp":  time.Now().Add(app_conf.SessionExpire()).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "appbase",
		"aud":  "appbase",
		"role": user.Role,
		"auth": user.IsAuth,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(srv_sec.JwtSecret))
}

// setSession attaches session and user info to the context.
// func setSession(c *gin.Context, user *app_models.Users) error {
// 	expire := time.Now().Add(app_conf.SessionExpire())
// 	if expire.IsZero() || expire.Before(time.Now()) {
// 		expire = time.Now().Add(24 * time.Hour)
// 	}

// 	locale := c.GetHeader("Accept-Language")
// 	if locale == "" {
// 		locale = "en-US"
// 	}

// 	session := app_models.Session{
// 		SessionID: srv_sec.UUID_String(),
// 		UserID:    user.ID,
// 		UserUUID:  user.UUID,
// 		Locale:    locale,
// 		Expire:    expire,
// 	}
// 	c.Set("appsession", session)
	
// 	return nil
// }

// setSecurityHeaders sets recommended security headers.
func setSecurityHeaders(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
}

// setCSRFCookie sets a CSRF token cookie.
func setCSRFCookie(c *gin.Context) {
	c.SetCookie("csrf_token", srv_sec.JwtSecret, app_conf.CookieExpire, "/", "", true, true)
}
