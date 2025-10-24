package app_db

import (
	"errors"
	"log"
	"time"

	"app/app_models"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user: something went wrong")
)

// CheckPassword compares a plaintext password with the hashed password
func CheckPassword(login_psw, user_psw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user_psw), []byte(login_psw))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func User_GetById(id any) (app_models.Users, error) {
	var userbyid app_models.Users
	err := AppDB.Where("id = ?", id).First(&userbyid)
	if err.Error != nil {
		log.Println("Error getting user by ID:", id, ErrUserNotFound)
		return userbyid, ErrUserNotFound
	}

	return userbyid, nil
}

func User_GetByUUID(uuid any) (app_models.Users, error) {
	var userbyuuid app_models.Users
	err := AppDB.
		Preload("AuthLevel").
		Preload("Org").
		Where("uuid = ?", uuid).First(&userbyuuid)
	if err.Error != nil {
		log.Println("Error getting user by UUID:", uuid, ErrUserNotFound)
		return userbyuuid, ErrUserNotFound
	}

	return userbyuuid, nil
}

func User_GetByUsername(username string) (app_models.Users, error) {
	var user app_models.Users
	err := AppDB.
		Preload("AuthLevel").
		Preload("Org").
		Where("username = ?", username).First(&user)
	if err.Error != nil {
		return user, ErrUserNotFound
	}
	return user, nil
}

func User_SetLastLogin(uuid any) error {
	// Update the last login time of the user
	err := AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("last_login", time.Now())
	if err.Error != nil {
		log.Println("Error updating last login:", err.Error)
		return err.Error
	}
	//log.Println("Last login updated successfully for user", uuid, time.Now())
	return nil
}

func Users_Count() int {
	var users_count int64
	AppDB.Model(&app_models.Users{}).Count(&users_count)
	return int(users_count)
}

func CreateNewUser(nu *app_models.Users) error {

	log.Println("Creating new user", nu.Username)
	res := *AppDB.Where("username", nu.Username).
		Attrs(app_models.Users{Username: nu.Username, Email: nu.Email, Password: nu.Password, AuthLevelID: nu.AuthLevelID, IsAuth: nu.IsAuth}).
		FirstOrCreate(&nu)

	if res.Error != nil {
		log.Println("Error creating user:", res.Error)
		return ErrUserNotFound
	}

	if res.RowsAffected == 0 { // if user already exists
		return errors.New("user already exists")
	}

	return nil
}

func Users_GetAll() ([]*app_models.Users, error) {
	var all_users []*app_models.Users
	err := *AppDB.Preload("AuthLevel").Preload("Org").Find(&all_users)
	if err.Error != nil {
		return all_users, ErrUserNotFound
	}
	return all_users, nil
}

func User_Delete(uuid any) error {
	// Delete user after updating updated_at
	err := AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("updated_at", AppDB.NowFunc())
	if err.Error != nil {
		return err.Error
	}
	// Delete user
	err = AppDB.Where("uuid = ?", uuid).Delete(&app_models.Users{})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func User_GetUsernameFromId(id any) (string, error) {
	var _user app_models.Users
	err := *AppDB.Where("id = ?", id).First(&_user)
	if err.Error != nil {
		return "", err.Error
	}
	return _user.Username, nil
}

func Users_GetAuth() ([]app_models.Users, error) {
	var auth_users []app_models.Users
	err := *AppDB.
		Preload("AuthLevel").
		Preload("Org").
		Where("is_auth = ?", true).Find(&auth_users)
	if err.Error != nil {
		return auth_users, err.Error
	}
	return auth_users, nil
}

func Users_GetUnAuth() ([]app_models.Users, error) {
	var unauth_users []app_models.Users
	err := *AppDB.
		Preload("AuthLevel").
		Preload("Org").
		Where("is_auth = ?", false).Find(&unauth_users)
	if err.Error != nil {
		return unauth_users, err.Error
	}
	return unauth_users, nil
}

func Users_GetDeleted() ([]app_models.Users, error) {
	var del_users []app_models.Users
	err := *AppDB.Unscoped().
		Preload("AuthLevel").
		Preload("Org").
		Where("deleted_at IS NOT NULL").Find(&del_users)
	if err.Error != nil {
		return del_users, err.Error
	}

	return del_users, nil
}

func Users_GetNew() ([]app_models.Users, error) {
	var new_users []app_models.Users
	err := AppDB.Where("created_at = updated_at and is_auth = false").Find(&new_users)
	if err.Error != nil {
		return new_users, err.Error
	}
	return new_users, nil
}

func GetAuthLevels() []app_models.AuthLevel {
	var auth_levels []app_models.AuthLevel
	err := AppDB.Find(&auth_levels)
	if err.Error != nil {
		return nil
	}
	return auth_levels
}
