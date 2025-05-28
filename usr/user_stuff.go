package users

import (
	"errors"
	"log"

	"app/app_db"
	"app/app_models"
)

func User_GetById(id any) (app_models.Users, error) {
	var userbyid app_models.Users
	err := app_db.AppDB.Where("id = ?", id).First(&userbyid)
	if err.Error != nil {
		log.Println("Error getting user by ID:", id, err.Error)
		return userbyid, err.Error
	}

	return userbyid, nil
}

func User_GetEmailById(userid any) (string, error) {
	var emailbyid app_models.Users
	err := app_db.AppDB.Select("email").First(&emailbyid, userid)
	if err.Error != nil {
		return "", err.Error
	}
	return emailbyid.Email, nil
}

func User_GetByEmail(email string) (app_models.Users, error) {
	var userbyemail app_models.Users
	err := app_db.AppDB.Where("email = ?", email).First(&userbyemail)
	if err.Error != nil {
		return userbyemail, err.Error
	}
	return userbyemail, nil
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
		return errors.New("error creating user")
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
		return all_users, err.Error
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

func user_SetNewPassword(uuid any, password string) error {
	err := app_db.AppDB.Model(&app_models.Users{}).Where("uuid = ?", uuid).Update("password", password)
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
	var role_user *app_models.Users
	err := app_db.AppDB.Where("id = ?", id).First(&role_user)
	if err.Error != nil {
		return "", err.Error
	}
	return role_user.Role, nil
}

func Users_GetAuth() ([]app_models.Users, int, error) {
	var auth_users *[]app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", true).Find(&auth_users)
	if err.Error != nil {
		return *auth_users, 0, err.Error
	}
	return *auth_users, len(*auth_users), nil
}

func Users_GetUnAuth() ([]app_models.Users, int, error) {
	var unauth_users *[]app_models.Users
	err := *app_db.AppDB.Where("is_auth = ?", false).Find(&unauth_users)
	if err.Error != nil {
		return *unauth_users, 0, err.Error
	}
	return *unauth_users, len(*unauth_users), nil
}

func Users_GetDeleted() ([]app_models.Users, int, error) {
	var del_users *[]app_models.Users
	err := *app_db.AppDB.Unscoped().Where("deleted_at IS NOT NULL").Find(&del_users)
	if err.Error != nil {
		return *del_users, 0, err.Error
	}

	return *del_users, len(*del_users), nil
}

func Users_GetNew() ([]app_models.Users, int, error) {
	var new_users []app_models.Users
	err := app_db.AppDB.Where("created_at = updated_at and is_auth = false").Find(&new_users)
	if err.Error != nil {
		return new_users, 0, err.Error
	}
	return new_users, len(new_users), nil
}
