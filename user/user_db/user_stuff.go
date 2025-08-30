package user_db

import (
	"errors"
	"log"
	"time"

	"app/app_db"
	"app/app_models"

	"golang.org/x/crypto/bcrypt"

)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user: something went wrong")
)

// CheckPassword compares a plaintext password with the hashed password
func CheckPassword(userid uint,password string) bool {
	var usr app_models.Users
	app_db.AppDB.First(&usr, "id = ?", userid)
	if usr.ID == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
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
	err := app_db.AppDB.Where("id = ?", id).First(&userbyid)
	if err.Error != nil {
		log.Println("Error getting user by ID:", id, ErrUserNotFound)
		return userbyid, ErrUserNotFound
	}

	return userbyid, nil
}

func User_GetEmailById(userid any) (string, error) {
	var emailbyid app_models.Users
	err := app_db.AppDB.Select("email").First(&emailbyid, userid)
	if err.Error != nil {
		return "", ErrUserNotFound
	}
	return emailbyid.Email, nil
}

func User_GetByEmail(email string) (app_models.Users, error) {
	var userbyemail app_models.Users
	err := app_db.AppDB.Where("email = ?", email).First(&userbyemail)
	if err.Error != nil {
		return userbyemail, ErrUserNotFound
	}
	return userbyemail, nil
}

func User_SetLastLogin(uuid uint) error {
	// Update the last login time of the user
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("last_login", time.Now())
	if err.Error != nil {
		log.Println("Error updating last login:", err.Error)
		return err.Error
	}
	//log.Println("Last login updated successfully for user", uuid, time.Now())
	return nil
}

func Users_Count() int {
	var users_count int64
	app_db.AppDB.Model(&app_models.Users{}).Count(&users_count)
	return int(users_count)
}

func CreateNewUser(nu *app_models.Users) error {

	log.Println("Creating new user", nu.Email)
	res := *app_db.AppDB.Where("email", nu.Email).
		Attrs(app_models.Users{Email: nu.Email, Password: nu.Password, Role: nu.Role, IsAuth: nu.IsAuth}).
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

func Users_GetAll() ([]app_models.Users, error) {
	var all_users []app_models.Users
	err := *app_db.AppDB.Find(&all_users)
	if err.Error != nil {
		return all_users, ErrUserNotFound
	}
	return all_users, nil
}

func User_Delete(uuid any) error {
	// Delete user after updating updated_at
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("updated_at", app_db.AppDB.NowFunc())
	if err.Error != nil {
		return err.Error
	}
	// Delete user
	err = app_db.AppDB.Where("uuid = ?", uuid).Delete(&app_models.Users{})
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func user_UpdateAuth(uuid any, isauth bool) error {
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("is_auth", isauth)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func user_UpdateRole(uuid any, role string) error {
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("role", role)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func user_UpdateOrg(uuid any, orgname string) error {
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("org", orgname)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func User_GetEmailFromId(id any) (string, error) {
	var user_email app_models.Users
	err := *app_db.AppDB.Where("id = ?", id).First(&user_email)
	if err.Error != nil {
		return "", err.Error
	}
	return user_email.Email, nil
}

func User_GetRoleFromId(id any) (string, error) {
	var role_user app_models.Users
	err := app_db.AppDB.Where("id = ?", id).First(&role_user)
	if err.Error != nil {
		return "", err.Error
	}
	return role_user.Role, nil
}

func Users_GetAuth() ([]app_models.Users, error) {
	var auth_users []app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", true).Find(&auth_users)
	if err.Error != nil {
		return auth_users, err.Error
	}
	return auth_users, nil
}

func Users_GetUnAuth() ([]app_models.Users, error) {
	var unauth_users []app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", false).Find(&unauth_users)
	if err.Error != nil {
		return unauth_users, err.Error
	}
	return unauth_users, nil
}

func Users_GetDeleted() ([]app_models.Users, error) {
	var del_users []app_models.Users
	err := *app_db.AppDB.Unscoped().Where("deleted_at IS NOT NULL").Find(&del_users)
	if err.Error != nil {
		return del_users, err.Error
	}

	return del_users, nil
}

func Users_GetNew() ([]app_models.Users, error) {
	var new_users []app_models.Users
	err := app_db.AppDB.Where("created_at = updated_at and is_auth = false").Find(&new_users)
	if err.Error != nil {
		return new_users, err.Error
	}
	return new_users, nil
}
