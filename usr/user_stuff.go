package users

import (
	"errors"
	"log"

	"app/app_db"
	"app/app_models"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user: something went wrong")
	dbUser = app_models.Users{}
	dbUsers = []app_models.Users{}
)




func User_GetById(id any) (app_models.Users, error) {
	//var userbyid app_models.Users
	err := app_db.AppDB.Where("id = ?", id).First(&dbUser)
	if err.Error != nil {
		log.Println("Error getting user by ID:", id, ErrUserNotFound)
		return dbUser, ErrUserNotFound
	}

	return dbUser, nil
}

func User_GetEmailById(userid any) (string, error) {
	//var emailbyid app_models.Users
	err := app_db.AppDB.Select("email").First(&dbUser, userid)
	if err.Error != nil {
		return "", ErrUserNotFound
	}
	return dbUser.Email, nil
}

func User_GetByEmail(email string) (app_models.Users, error) {
	//var userbyemail app_models.Users
	err := app_db.AppDB.Where("email = ?", email).First(&dbUser)
	if err.Error != nil {
		return dbUser, ErrUserNotFound
	}
	return dbUser, nil
}

func User_SetLastLogin(uuid int) error {
	// Update the last login time of the user
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("last_login", app_db.AppDB.NowFunc())
	if err.Error != nil {
		log.Println("Error updating last login:", err.Error)
		return err.Error
	}
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
	//var all_users []app_models.Users
	err := *app_db.AppDB.Find(&dbUsers)
	if err.Error != nil {
		return dbUsers, ErrUserNotFound
	}
	return dbUsers, nil
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

func user_SetNewPassword(uuid any, password string) error {
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("password", password)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func User_GetEmailFromId(id any) (string, error) {
	//var user_email app_models.Users
	err := *app_db.AppDB.Where("id = ?", id).First(&dbUser)
	if err.Error != nil {
		return "", err.Error
	}
	return dbUser.Email, nil
}

func User_GetRoleFromId(id any) (string, error) {
	// var role_user app_models.Users
	err := app_db.AppDB.Where("id = ?", id).First(&dbUser)
	if err.Error != nil {
		return "", err.Error
	}
	return dbUser.Role, nil
}

func Users_GetAuth() ([]app_models.Users, int, error) {
	//var auth_users []app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", true).Find(&dbUsers)
	if err.Error != nil {
		return dbUsers, 0, err.Error
	}
	return dbUsers, len(dbUsers), nil
}

func Users_GetUnAuth() ([]app_models.Users, int, error) {
	//var unauth_users []app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", false).Find(&dbUsers)
	if err.Error != nil {
		return dbUsers, 0, err.Error
	}
	return dbUsers, len(dbUsers), nil
}

func Users_GetDeleted() ([]app_models.Users, int, error) {
	//var del_users []app_models.Users
	err := *app_db.AppDB.Unscoped().Where("deleted_at IS NOT NULL").Find(&dbUsers)
	if err.Error != nil {
		return dbUsers, 0, err.Error
	}

	return dbUsers, len(dbUsers), nil
}

func Users_GetNew() ([]app_models.Users, int, error) {
	//var new_users []app_models.Users
	err := app_db.AppDB.Where("created_at = updated_at and is_auth = false").Find(&dbUsers)
	if err.Error != nil {
		return dbUsers, 0, err.Error
	}
	return dbUsers, len(dbUsers), nil
}
