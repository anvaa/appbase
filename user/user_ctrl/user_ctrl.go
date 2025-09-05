package user_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"app/app_models"
	
	"server/middleware"
	"user/user_conf"
	"user/user_sec"

	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

var (
	errmsg = "user or password invalid"

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

func SignUp(c *gin.Context) {

	ip := c.ClientIP()
	if !canAttemptLogin(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Please wait before trying again."})
		return
	}

	var body struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Password2 string `json:"password2"`
		Orgname   string `json:"orgname"`
		Count     int    `json:"count"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read body"})
		return
	}

	email, password, password2, orgname := strings.TrimSpace(body.Email), strings.TrimSpace(body.Password), strings.TrimSpace(body.Password2), strings.TrimSpace(body.Orgname)

	if err := user_sec.IsValidEmail(email); err != nil {
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

	user, err := app_db.User_GetByEmail(email)
	if err != nil {
		log.Println("User not found >> New User")
	}

	if user.ID > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	role, isauth := 4, false
	user = app_models.Users{
		Email:       email,
		AuthLevelID: role,
		IsAuth:      isauth,
		Note:        "Nil",
		Org:         &[]app_models.Org{{Name: orgname}},
	}

	hash, err := hashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}
	user.Password = hash

	if err := app_db.CreateNewUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	url := "/login"
	if body.Count == 121209 {
		url = "/v/users"
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "url": url})
}

func Login(c *gin.Context) {
	// Rate limit: allow login attempt only every 3 seconds per IP
	ip := c.ClientIP()
	if !canAttemptLogin(ip) {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Please wait before trying again."})
		return
	}

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read"})
		return
	}

	email, password := strings.TrimSpace(body.Email), strings.TrimSpace(body.Password)
	url := "/"

	if err := user_sec.IsValidEmail(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	if err := user_sec.IsValidPassword(password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	user, err := app_db.User_GetByEmail(email)
	if err != nil || !user.IsAuth || user.Email == "" || !app_db.CheckPassword(password, user.Password) {
		log.Println("Login failed:", email)
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	err = middleware.SetJWT(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set JWT"})
		return
	}

	err = app_db.User_SetLastLogin(user.UUID)
	if err != nil {
		log.Println("Error setting last login:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set last login"})
		return
	}

	url = app_conf.GetString("start_url")
	c.JSON(http.StatusOK, gin.H{"url": url, "message": "success"})
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
	var body struct {
		Uuid string `json:"uuid"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
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
		Uuid string `json:"uuid"`
		Auth bool   `json:"auth"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	if body.Auth {
		body.Auth = false
	} else {
		body.Auth = true
	}

	// update user
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
		Uuid    string `json:"uuid"`
		AuthLev int    `json:"authlevel"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
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
		Uuid    string `json:"uuid"`
		Name    string `json:"name"`
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

	user.Org = &[]app_models.Org{{Name: body.Name}}

	if err := app_db.AppDB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_SetNewPassword(c *gin.Context) {
	var body struct {
		Uuid string `json:"uuid"`
		Psw1 string `json:"psw1"`
		Psw2 string `json:"psw2"`
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
	err := app_db.AppDB.Where("uuid = ?", body.Uuid).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not found"})
		return
	}

	// generate new hash from password
	hash, err := hashPassword(body.Psw1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}

	// update user
	user.Password = string(hash)

	// Increment TokenVersion to invalidate all existing JWTs
	user.TokenVersion++
	err = app_db.AppDB.Model(&user).Updates(map[string]interface{}{"password": user.Password, "token_version": user.TokenVersion}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update password"})
		return
	}

	// Invalidate the auth cookie so the user must re-login
	c.SetCookie(user_conf.CookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
