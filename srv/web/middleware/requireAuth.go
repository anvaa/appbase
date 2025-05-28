package middleware

import (
	"srv/srv_sec"
	"app/app_conf"
	"usr"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userKey string = "user"

func RequireAuth(c *gin.Context) {
	// Get the JWT string from the header
	tokenString, err := c.Cookie(app_conf.CookieName)
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

	user, exists := users.User_GetById(claims["sub"])
	if exists != nil {
		redirectToLogin(c)
		return
	}

	// Check if the user is authenticated
	if !user.IsAuth {
		redirectToLogin(c)
		return
	}

	// Attach the user to the context
	c.Set(userKey, user)
	c.Next()
}

func redirectToLogin(c *gin.Context) {
	c.SetCookie(app_conf.CookieName, "", -1, "/", "", false, true)
	c.Redirect(302, "/login")
	c.Abort()
}
