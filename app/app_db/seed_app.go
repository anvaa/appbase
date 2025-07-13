package app_db

import (
	"app/app_models"
	"fmt"
	"gorm.io/gorm"
	
)

func seedProject(db *gorm.DB) {
	// Check if the number of projects is > 0
	var count int64
	db.Model(&app_models.Project{}).Count(&count)
	if count > 0 {
		fmt.Println("Projects already seeded")
		return
	}

	projects := []app_models.Project{
		{Name: "Project Alpha", Note: "First project", StasubID : 1, TypsubID: 1},
		{Name: "Project Beta", Note: "", StasubID: 2, TypsubID: 2},
	}

	for i := range projects {
		db.Create(&projects[i])
	}

	fmt.Println("Projects Seeded:", projects[0].Name, projects[1].Name)
}