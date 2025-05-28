package users

import (
	"app/app_conf"
	"app/app_models"

	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	errmsg = "user or password invalid"
)

func SignUp(c *gin.Context) {
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

	if err := IsValidEmail(email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if password != password2 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Passwords do not match"})
		return
	}

	if err := IsValidPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	if err := IsValidOrg(orgname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err := User_GetByEmail(email)
	if err != nil {
		log.Println("User not found >> New User")
	}

	if user.ID > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
		return
	}

	role, isauth := "user", false

	user = app_models.Users{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		IsAuth:   isauth,
		Note:     "Nil",
	}

	if err := CreateNewUser(&user); err != nil {
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

	if err := IsValidEmail(email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	if err := IsValidPassword(password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	user, err := User_GetByEmail(email)
	if err != nil || !user.IsAuth || user.Email == "" || !CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": errmsg, "url": url})
		return
	}

	c, err = SetJWT(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to set JWT"})
		return
	}

	url = app_conf.GetString("start_url")
	c.JSON(http.StatusOK, gin.H{"url": url, "message": "success"})
}

func GetAllUsers(c *gin.Context) {
	users, err := Users_GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := User_GetById(id)
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

	if err := User_Delete(body.Uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func User_UpdateAuth(c *gin.Context) {
	var body struct {
		Uuid int  `json:"uuid"`
		Auth bool `json:"auth"`
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

	if err := user_UpdateAuth(body.Uuid, body.Auth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update auth"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_UpdateRole(c *gin.Context) {
	var body struct {
		Uuid string `json:"uuid"`
		Role string `json:"role"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	if err := user_UpdateRole(body.Uuid, body.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update role"})
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

	if err := IsValidPassword(body.Psw1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Password is not valid"})
		return
	}

	hashedPassword, err := HashPassword(body.Psw1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	if err := user_SetNewPassword(body.Uuid, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func User_UpdOrg(c *gin.Context) {
	var body struct {
		Uuid    string `json:"uuid"`
		Orgname string `json:"orgname"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	orgname := strings.TrimSpace(body.Orgname)
	if err := user_UpdateOrg(body.Uuid, orgname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update orgname time"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
