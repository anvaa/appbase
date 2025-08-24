package app_db

import (
	"app/app_models"
	
)

func Sta_GetStaIDByTitle(title string) (int, error) {
	var sta app_models.Status
	if err := AppDB.Where("title = ?", title).First(&sta).Error; err != nil {
		return 0, err
	}
	return sta.ID, nil
}

func Sta_GetAllStasub(title string) ([]app_models.Stasub, error) {

	sta_id, err := Sta_GetStaIDByTitle(title)
	if err != nil || sta_id == 0 {
		return nil, err
	}

	var stasub []app_models.Stasub
	if err = AppDB.Where("status_id = ?", sta_id).Find(&stasub).Error; err != nil {
		return nil, err
	}

	return stasub, nil
}

func Sta_GetNameById(id any) string {
	var sta app_models.Stasub
	if err := AppDB.Where("id = ?", id).First(&sta).Error; err != nil {
		return "nil"
	}

	return sta.Name
}

func Sta_GetStaSubIDByType(val string) []app_models.Stasub {
	var sta []app_models.Stasub
	if err := AppDB.Where("type = ?", val).Find(&sta).Error; err != nil {
		return nil
	}

	return sta
}