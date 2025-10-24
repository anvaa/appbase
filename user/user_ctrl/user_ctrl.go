package user_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"app/app_models"
	"fmt"
	"log"
	"net/http"

	"strings"
	"sync"
	"time"

	"server/middleware"
	"server/srv_conf"
	"user/user_conf"
	"user/user_sec"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	errmsg         = "user or password invalid"
	loginAttempts  = make(map[string]time.Time)
	loginMu        sync.Mutex
	loginRateLimit = user_conf.LoginRateLimit()
)

func canAttemptLogin(ip string) bool {
	loginMu.Lock()
	defer loginMu.Unlock()
	last, exists := loginAttempts[ip]
	if !exists || time.Since(last) > loginRateLimit {
		loginAttempts[ip] = time.Now()
		return true
	}
	return false
}

func Health(c *gin.Context) {
	// json { "status": "ok", "timestamp": "2025-01-01T00:00:00Z" }
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func SignUp(c *gin.Context) {
	ip := c.ClientIP()
	if !canAttemptLogin(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Please wait before trying again."})
		return
	}

	var body struct {
		Username, Email, Password, Password2, Orgname string
		Count                                         int
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read body"})
		return
	}

	username, email, password, password2, orgname := strings.TrimSpace(body.Username),
		strings.TrimSpace(body.Email),
		strings.TrimSpace(body.Password),
		strings.TrimSpace(body.Password2),
		strings.TrimSpace(body.Orgname)
	if err := user_sec.IsValidUsername(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if password != password2 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Passwords do not match"})
		return
	}

	if err := user_sec.IsValidPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}
	if err := user_sec.IsValidOrg(orgname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err := app_db.User_GetByUsername(username)
	if err != nil {
		log.Println("User not found >> New User")
	}
	if user.ID > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	hash, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	newUser := app_models.Users{
		Username:    username,
		Email:       email,
		AuthLevelID: 4,
		IsAuth:      false,
		Note:        "NULL",
		Org:         []*app_models.Org{{Name: orgname}},
		Password:    hash,
	}

	if err := app_db.CreateNewUser(&newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	url := "/login"
	if body.Count == 121209 {
		url = "/v/users"
	}

	log.Println("New user created:", username, "Redirect:", url)
	c.JSON(http.StatusOK, gin.H{"message": "success", "redirect": url})
}

func Login(c *gin.Context) {
	ip := c.ClientIP()
	if !canAttemptLogin(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Please wait before trying again."})
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	username := strings.TrimSpace(body.Username)
	password := strings.TrimSpace(body.Password)
	// fmt.Println("Login attempt for:", username)

	if err := user_sec.IsValidUsername(username); err != nil || user_sec.IsValidPassword(password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errmsg})
		return
	}

	user, err := app_db.User_GetByUsername(username)

	if err != nil || user.Username == "" || !user.IsAuth || !app_db.CheckPassword(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errmsg})
		return
	}

	// fmt.Println("User found:", user.Username, "AuthLevel:", user.AuthLevel.Name)

	tokenString, err := middleware.SetJWT(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to set authentication"})
		return
	}

	if err := app_db.User_SetLastLogin(user.UUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update last login"})
		return
	}

	url := app_conf.BaseURL()
	if url == "" {
		url = "/"
	}
	fmt.Println("User logged in:", username, "Redirect:", url)

	c.JSON(http.StatusOK, gin.H{
		"token":        tokenString,
		"session":      tokenString,
		"access_token": tokenString,
		"redirect":     url,
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
			"roles":    []string{user.AuthLevel.Name},
			"profile": gin.H{
				"uuid":       user.UUID,
				"auth_level": user.AuthLevel.Name,
				"orgs":       user.Org,
				"note":       user.Note,
			},
		},
	})

}

func Logout(c *gin.Context) {
	secure := srv_conf.UseTLS()
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", secure, true)
	// c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	c.Redirect(http.StatusFound, "/login")
}

func GetAllUsers(c *gin.Context) {
	users, err := app_db.Users_GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := app_db.User_GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func User_DeleteUser(c *gin.Context) {
	var body struct{ Uuid string }
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	user, err := app_db.User_GetByUUID(body.Uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}
	if user.AuthLevel.Level >= 40 {
		c.JSON(http.StatusForbidden, gin.H{"message": "not allowed to delete admin level user"})
		return
	}

	curUser := c.Keys[user_conf.UserKey].(app_models.Users)
	if curUser.UUID == user.UUID {
		c.JSON(http.StatusForbidden, gin.H{"message": "not allowed to delete yourself"})
		return
	}

	if err := app_db.User_Delete(body.Uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func User_UpdateAuth(c *gin.Context) {
	var body struct {
		Uuid string
		Auth bool
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}
	body.Auth = !body.Auth

	if err := app_db.AppDB.Model(&app_models.Users{}).
		Where("uuid = ?", body.Uuid).
		Update("is_auth", body.Auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_UpdAuthLevel(c *gin.Context) {
	var body struct {
		Uuid    string
		AuthLev int
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	user, err := app_db.User_GetByUUID(body.Uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}

	curUser := c.Keys[user_conf.UserKey].(app_models.Users)
	if curUser.UUID == user.UUID {
		c.JSON(http.StatusForbidden, gin.H{"message": "not allowed to edit yourself"})
		return
	}

	if err := app_db.AppDB.Model(&app_models.Users{}).
		Where("uuid = ?", body.Uuid).
		Update("auth_level_id", body.AuthLev).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_UpdateOrg(c *gin.Context) {
	var body struct {
		Uuid, Name string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	user, err := app_db.User_GetByUUID(body.Uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
		return
	}
	user.Org = []*app_models.Org{{Name: body.Name}}

	if err := app_db.AppDB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_SetNewPassword(c *gin.Context) {
	var body struct {
		Uuid, Psw1, Psw2 string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}
	if err := user_sec.IsValidPassword(body.Psw1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password is not valid"})
		return
	}

	var user app_models.Users
	if err := app_db.AppDB.Where("uuid = ?", body.Uuid).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		return
	}

	hash, err := hashPassword(body.Psw1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}
	user.Password = hash
	user.TokenVersion++

	if err := app_db.AppDB.Model(&user).Updates(map[string]interface{}{
		"password":      user.Password,
		"token_version": user.TokenVersion,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update password"})
		return
	}

	secure := srv_conf.UseTLS()
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
