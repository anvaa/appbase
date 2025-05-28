package app_db

import (
	"app/app_models"
	
)



func Sta_GetNameById(id any) string {
	var sta app_models.Status
	if err := AppDB.Where("id = ?", id).First(&sta).Error; err != nil {
		return "nil"
	}

	return sta.Name
}

func Sta_GetStatusIdByType(val string) int {
	var sta app_models.Status
	if err := AppDB.Where("type = ?", val).First(&sta).Error; err != nil {
		return 0
	}

	return sta.ID
}

func Sta_HistoryDelete(itmid any) error {
	// Delete the status history of an item
	return AppDB.Where("itmid = ?", itmid).Delete(&app_models.StatusHistory{}).Error
}