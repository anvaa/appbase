package users

import (
	"errors"
	"regexp"
	"srv/srv_sec"

	"golang.org/x/crypto/bcrypt"

	"app/app_conf"
	"app/app_models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ErrSqlInText = regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|--|;)`)
	ErrIsAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	ErrIsAlphanumericWithSpace = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	ErrIsAlphanumericWithUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	ErrIsAlphanumericWithUnderscoreAndSpace = regexp.MustCompile(`^[a-zA-Z0-9_ ]+$`)
	ErrInvalidEmail = errors.New("not a valid email")
	ErrEmptyPassword = errors.New("password can't be empty")
	ErrShortPassword = errors.New("password must be at least 8 characters")
	ErrLongPassword  = errors.New("password must be less than 255 characters")
)

func IsValidEmail(email string) error {
	
	// Check if the email is empty
	if email == "" {
		return ErrInvalidEmail
	}
	// Check if the email is too long
	if len(email) > 255 {
		return ErrInvalidEmail
	}

	// Check for SQL injection
	if ErrSqlInText.MatchString(email) {
		return errors.New("sql injection detected")
	}

	// Check if the email matches the regex pattern
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func IsValidPassword(psw string) error {

	// Check for SQL injection
	if ErrSqlInText.MatchString(psw) {
		return errors.New("sql injection detected")
	}

	switch {
	case psw == "":
		return ErrEmptyPassword
	case len(psw) < 8:
		return ErrShortPassword
	case len(psw) > 255:
		return ErrLongPassword
	default:
		return nil
	}
}

func IsValidOrg(org string) error {

	if org == "" {
		return errors.New("organization can't be empty")
	}

	// Check for SQL injection
	if ErrSqlInText.MatchString(org) {
		return errors.New("sql injection detected")
	}

	if ErrIsAlphanumericWithUnderscoreAndSpace.MatchString(org) {
		return nil
	}

	if len(org) > 255 {
		return errors.New("organization must be less than 255 characters")
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUuid() int {
	return int(uuid.New().ID())
}

func SetJWT(c *gin.Context, user *app_models.Users) (*gin.Context, error) {
	// Generate a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Second * time.Duration(app_conf.CookieExpire)).Unix(),
		"iat": time.Now().Unix(),
		"email": user.Email,
		"role": user.Role,
		"isauth": user.IsAuth,
	})

	tokenString, err := token.SignedString([]byte(srv_sec.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate"})
		return nil, err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(app_conf.CookieName, tokenString, app_conf.CookieExpire, "/", "", true, true)

	return c, nil
}