package middleware

import (
	
	"app/app_db"
	"server/srv_sec"
	"user/user_conf"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	userKey  = user_conf.UserKey
	UserUUID uint
	UserID   uint
)

func IsAuth(c *gin.Context) {
	// Get the JWT string from the header
	tokenString, err := c.Cookie(user_conf.CookieName)
	if err != nil {
		redirectToLogin(c)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(srv_sec.JwtSecret), nil
	},
	)

	if err != nil {
		redirectToLogin(c)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		OnErr(c)
		return
	}

	user, exists := app_db.User_GetByUUID(claims["sub"])
	if exists != nil {
		redirectToLogin(c)
		return
	}

	// Check if the user is authenticated
	if !user.IsAuth {
		redirectToLogin(c)
		return
	}

	// Check token_version for JWT invalidation
	tokenVersion, ok := claims["token_version"].(float64)
	if !ok || int(tokenVersion) != user.TokenVersion {
		// Token is invalid due to version mismatch
		redirectToLogin(c)
		return
	}

	// Attach the user to the context
	c.Set(userKey, user)
	UserID = user.ID
	UserUUID = user.UUID
	c.Next()
}

func redirectToLogin(c *gin.Context) {
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", false, true)
	c.Redirect(302, "/login")
	c.Abort()
}
