package app_db

import (
	"app/app_models"
	"fmt"

)

func GetAllProjects() []app_models.Project {
	var projects []app_models.Project
	if err := AppDB.
		Preload("Persons").
		Preload("Notes").
		//Preload("Phones").
		//Preload("Addresses").
		Preload("Emails").

		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").

		Preload("Stasub").
		Preload("Typsub").
		
		Order("Stasub_id ASC").
		Find(&projects).Error; err != nil {
		fmt.Println("Error fetching projects:", err)
		return nil
	}
	return projects
}

func GetProjectByUUID(uuid any) (*app_models.Project, error) {
	var project app_models.Project
	if err := AppDB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Preload("Stasub").
		Preload("Typsub").
		
		Preload("Persons.Typsub").
		Preload("Persons.Notes").
		Preload("Persons.Emails").

		Preload("Notes.Typsub").
		Preload("Emails.Typsub").
		Preload("Emails.Stasub").


		First(&project, "uuid = ?", uuid).Error; err != nil {
		return nil, fmt.Errorf("error fetching project with UUID %d: %w", uuid, err)
	}
	return &project, nil
}
