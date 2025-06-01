package app_db

import (
	"app/app_models"
	
)



func Sta_GetNameById(id any) string {
	var sta app_models.StaSub
	if err := AppDB.Where("id = ?", id).First(&sta).Error; err != nil {
		return "nil"
	}

	return sta.Name
}

func Sta_GetStaSubUUIDByType(val string) []app_models.StaSub {
	var sta []app_models.StaSub
	if err := AppDB.Where("type = ?", val).Find(&sta).Error; err != nil {
		return nil
	}

	return sta
}