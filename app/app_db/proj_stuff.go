package app_db

import (
	"app/app_models"
	"fmt"

)

func GetAllProjects() []app_models.Project {
	var projects []app_models.Project
	if err := AppDB.
		//Preload("Persons").
		//Preload("Phones").
		//Preload("Addresses").
		//Preload("Emails").
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
		//Preload("Persons").
		//Preload("Phones").
		//Preload("Addresses").
		//Preload("Emails").
		First(&project, "uuid = ?", uuid).Error; err != nil {
		return nil, fmt.Errorf("error fetching project with UUID %d: %w", uuid, err)
	}
	return &project, nil
}
