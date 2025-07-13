package middleware

import (
	"encoding/base64"
	"net/http"

	"srv/srv_conf"
	"srv/srv_sec"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFProtection sets up session middleware for Gin
func CSRFProtection() gin.HandlerFunc {
	store := cookie.NewStore([]byte(srv_sec.GetCSRFSecret()))
	return sessions.Sessions("appsession", store)
}

// CSRF returns a Gin middleware that wraps gorilla/csrf protection
func CSRF() gin.HandlerFunc {
	// Get the CSRF secret (base64 encoded)
	secretStr := srv_sec.GetCSRFSecret()

	// Decode the base64 secret
	csrfKey, err := base64.URLEncoding.DecodeString(secretStr)
	if err != nil {
		// Fallback: use the string directly and ensure 32 bytes
		csrfKey = []byte(secretStr)
	}

	// Ensure key is exactly 32 bytes for security
	if len(csrfKey) > 32 {
		csrfKey = csrfKey[:32]
	} else if len(csrfKey) < 32 {
		// Pad with zeros if too short (not ideal, but fallback)
		temp := make([]byte, 32)
		copy(temp, csrfKey)
		csrfKey = temp
	}

	// Configure CSRF protection with environment-specific secure options
	isSecure := !srv_conf.IsGinModDebug() // Only use secure cookies in production

	csrfProtection := csrf.Protect(
		csrfKey,
		csrf.Secure(isSecure),                  // Use secure cookies in production only
		csrf.HttpOnly(true),                    // Prevent XSS attacks
		csrf.SameSite(csrf.SameSiteStrictMode), // Strict same-site policy
		csrf.Path("/"),                         // Apply to entire site
		csrf.Domain(""),                        // Let browser determine domain
		csrf.MaxAge(3600),                      // Token expires in 1 hour
		csrf.RequestHeader("X-CSRF-Token"),     // Custom header name
		csrf.FieldName("csrf_token"),           // Form field name
		csrf.ErrorHandler(http.HandlerFunc(csrfErrorHandler)),
	)

	return func(c *gin.Context) {
		// Skip CSRF for certain content types (like API endpoints expecting JSON)
		if shouldSkipCSRF(c) {
			c.Next()
			return
		}

		// Create a dummy handler to wrap
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})

		// Apply CSRF protection
		csrfProtection(handler).ServeHTTP(c.Writer, c.Request)
	}
}

// CSRFToken returns the CSRF token for the current request
func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Make CSRF token available to templates and responses
		token := csrf.Token(c.Request)
		c.Set("csrf_token", token)
		c.Header("X-CSRF-Token", token)
		c.Next()
	}
}

// shouldSkipCSRF determines if CSRF protection should be skipped for this request
func shouldSkipCSRF(c *gin.Context) bool {
	// Skip CSRF for API endpoints that use other authentication methods
	contentType := c.GetHeader("Content-Type")

	// Skip for API requests with JSON content type
	if contentType == "application/json" {
		return true
	}

	// Skip for preflight OPTIONS requests
	if c.Request.Method == "OPTIONS" {
		return true
	}

	// Check if route is explicitly marked to skip CSRF
	if skip, exists := c.Get("csrf_skip"); exists && skip.(bool) {
		return true
	}

	return false
}

// CSRFErrorHandler handles CSRF validation errors
func csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Set appropriate headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	// Return JSON error response
	errorResponse := `{
		"error": "CSRF token validation failed",
		"message": "Invalid or missing CSRF token",
		"code": 403
	}`
	w.Write([]byte(errorResponse))
}

// SkipCSRF allows certain routes to bypass CSRF protection
func SkipCSRF(paths ...string) gin.HandlerFunc {
	skipMap := make(map[string]bool)
	for _, path := range paths {
		skipMap[path] = true
	}

	return func(c *gin.Context) {
		if skipMap[c.Request.URL.Path] {
			c.Set("csrf_skip", true)
		}
		c.Next()
	}
}

// CSRFTokenResponse returns the CSRF token as JSON response
func CSRFTokenResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.JSON(http.StatusOK, gin.H{
			"csrf_token": token,
		})
	}
}
