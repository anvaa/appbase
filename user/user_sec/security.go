package user_sec

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex                              = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	ErrSqlInText                            = regexp.MustCompile(`(?i)(SELECT|WHERE|INSERT|UPDATE|DELETE|DROP|--|;)`)
	ErrIsAlphanumeric                       = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	ErrIsAlphanumericWithSpace              = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)
	ErrIsAlphanumericWithUnderscore         = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	ErrIsAlphanumericWithUnderscoreAndSpace = regexp.MustCompile(`^[a-zA-Z0-9_ ]+$`)
	ErrInvalidUsername                      = errors.New("not a valid username")
	ErrInvalidEmail                         = errors.New("not a valid email")
	ErrEmptyPassword                        = errors.New("password can't be empty")
	ErrShortPassword                        = errors.New("password must be at least 8 characters")
	ErrLongPassword                         = errors.New("password must be less than 255 characters")
)

func IsValidUsername(username string) error {

	// Check if the username is empty
	if username == "" {
		return ErrInvalidUsername
	}
	// Check if the username is too long
	if len(username) > 255 {
		fmt.Println("Username is too long:", username)
		return ErrInvalidUsername
	}

	// Check for SQL injection
	if ErrSqlInText.MatchString(username) {
		fmt.Println("SQL Injection detected in username:", username)
		return errors.New("sql injection detected")
	}

	// If the username contains "@" and ".", validate it as an email
	if strings.Contains(username, "@") && strings.Contains(username, ".") {
		fmt.Println("Username is being validated as an email:", username)
		// Check if the username matches the email regex pattern
		if !emailRegex.MatchString(username) {
			fmt.Println("Email format is invalid:", username)
			return ErrInvalidEmail
		}
		return nil
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
