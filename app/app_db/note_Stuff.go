package app_db

import (
	"app/app_models"
	"fmt"

)

func GetAllNotes(proj_id, pers_id string) ([]app_models.Notes, error) {

	if proj_id == "" && pers_id == "" {
		return nil, fmt.Errorf("either project_id or person_id must be provided")
	}

	var seek_val string
	if proj_id != "" {
		seek_val = fmt.Sprintf("project_id = %s", proj_id)
	} 
	
	if pers_id != "" {
		seek_val = fmt.Sprintf("person_id = %s", pers_id)
	}

	var notes []app_models.Notes
	err := AppDB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Preload("Typsub").
		Order("updated_at DESC").
		Find(&notes, seek_val).Error
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