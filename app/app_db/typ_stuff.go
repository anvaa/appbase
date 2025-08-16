package app_db

import (
	"app/app_models"
	
)

func Typ_GetAllTypsub(typ int) []app_models.Typsub {
	var typsub []app_models.Typsub
	if err := AppDB.Where("type_id = ?", typ).Find(&typsub).Error; err != nil {
		return nil
	}

	return typsub
}

func Typ_GetTypSubIDByType(typ string) []app_models.Typsub {
	var typsub []app_models.Typsub
	if err := AppDB.Where("type = ?", typ).Find(&typsub).Error; err != nil {
		return nil
	}

	return typsub
}