package app_db

import (
	"app/app_models"
	"fmt"

)

func GetAllNotes(projid string, sort int) ([]app_models.Notes, error) {
	var notes []app_models.Notes
	err := AppDB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Preload("Typsub").
		Order("updated_at DESC").
		Find(&notes, "project_id = ?", projid).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get notes: %w", err)
	}
	return notes, nil
}

func GetNoteByUUID(uuid string) (*app_models.Notes, error) {
	var note app_models.Notes
	err := AppDB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Preload("Typsub").
		First(&note, "uuid = ?", uuid).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}
	return &note, nil
}