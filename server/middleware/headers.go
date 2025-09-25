package middleware

import (
	"net/http"
	"time"
	"user/user_conf"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"app/app_models"
	"server/srv_sec"
)

// SetJWT generates a JWT, sets it as a cookie, and attaches session/user info to the context.
func SetJWT(c *gin.Context, user *app_models.Users) (string, error) {
	tokenString, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return "", err
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(user_conf.CookieName, tokenString, user_conf.CookieExpire, "/", "", true, true)

	setSecurityHeaders(c)
	setCSRFCookie(c)

	return tokenString, nil
}

// generateJWT creates a signed JWT string for the given user.
func generateJWT(user *app_models.Users) (string, error) {
	claims := jwt.MapClaims{
		"sub":           user.UUID,
		"exp":           time.Now().Add(user_conf.SessionExpire()).Unix(),
		"iat":           time.Now().Unix(),
		"iss":           "appbase",
		"aud":           "appbase",
		"role":          user.AuthLevelID,
		"auth":          user.IsAuth,
		"token_version": user.TokenVersion,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(srv_sec.JwtSecret))
}

// setSecurityHeaders sets recommended security headers.
func setSecurityHeaders(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
}

// setCSRFCookie sets a CSRF token cookie.
func setCSRFCookie(c *gin.Context) {
	c.SetCookie("csrf_token", srv_sec.JwtSecret, user_conf.CookieExpire, "/", "", true, true)
}
