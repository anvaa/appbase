package app_db

import (
	"app/app_models"
	
)

func Typ_GetAllTypIDByTitle(title string) (uint, error) {
	var typ app_models.Type
	if err := AppDB.Where("title = ?", title).First(&typ).Error; err != nil {
		return 0, err
	}
	return typ.ID, nil
}

func Typ_GetAllTypsub(title string) ([]app_models.Typsub, error) {

	typ_id, err := Typ_GetAllTypIDByTitle(title)
	if err != nil {
		return nil, err
	}

	var typsub []app_models.Typsub
	if err = AppDB.Where("type_id = ?", typ_id).Find(&typsub).Error; err != nil {
		return nil, err
	}

	return typsub, nil
}

func Typ_GetTypSubIDByType(typ string) []app_models.Typsub {
	var typsub []app_models.Typsub
	if err := AppDB.Where("type = ?", typ).Find(&typsub).Error; err != nil {
		return nil
	}

	return typsub
}