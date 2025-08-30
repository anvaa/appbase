package app_db

import (
	"app/app_models"
)

func User_GetByUUID(uuid any) (app_models.Users, error) {
	var user app_models.Users
	result := AppDB.First(&user, "uuid = ?", uuid)
	if result.Error != nil {
		return app_models.Users{}, result.Error
	}
	return user, nil
}