package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"app/app_conf"
	"app/app_models"
	"server/srv_conf"
	"server/srv_sec"
	"user/user_conf"
)

// SetJWT generates a JWT, sets it as a cookie, and attaches session/user info to the context.
func SetJWT(c *gin.Context, user *app_models.Users) (string, error) {
	tokenString, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return "", err
	}

	c.SetSameSite(http.SameSiteStrictMode)

	// Secure: true ensures the cookie is only sent over HTTPS
	// HttpOnly: true prevents JavaScript access to the cookie
	c.SetCookie(user_conf.CookieName, tokenString, user_conf.CookieExpire, "/", "", srv_conf.UseTLS(), true)

	setSecurityHeaders(c)

	return tokenString, nil
}

// generateJWT creates a signed JWT string for the given user.
func generateJWT(user *app_models.Users) (string, error) {
	app_name := strings.Replace(app_conf.AppInfo().AppName, " ", "_", -1)
	app_name = strings.TrimSpace(app_name)
	claims := jwt.MapClaims{
		"sub":           user.UUID,
		"exp":           time.Now().Add(user_conf.SessionExpire()).Unix(),
		"iat":           time.Now().Unix(),
		"iss":           app_name,
		"aud":           app_name,
		"role":          user.AuthLevelID,
		"auth":          user.IsAuth,
		"token_version": user.TokenVersion,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(srv_sec.JWTSecret()))
}

// setSecurityHeaders sets recommended security headers.
func setSecurityHeaders(c *gin.Context) {
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self'")
	c.Header("Referrer-Policy", "no-referrer")
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

}
