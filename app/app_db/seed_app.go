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
		{Name: "Project Alpha", StasubID: 1, TypsubID: 1, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
		{Name: "Project Beta", StasubID: 2, TypsubID: 2, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
	}

	for i := range projects {
		db.Create(&projects[i])
	}

	fmt.Println("Projects Seeded:")
}

func seedNotes(db *gorm.DB) {
	// Check if the number of notes is > 0
	var count int64
	db.Model(&app_models.Notes{}).Count(&count)
	if count > 0 {
		fmt.Println("Notes already seeded")
		return
	}

	notes := []app_models.Notes{
		{Content: "Note 1", ProjectID: 1, TypsubID: 7},
		{Content: "Note 2", ProjectID: 1, TypsubID: 8},
		{Content: "Note 3", ProjectID: 1, TypsubID: 9},
		{Content: "Note 4", ProjectID: 2, TypsubID: 7},
		{Content: "Note 5", ProjectID: 2, TypsubID: 8},
	}

	for i := range notes {
		db.Create(&notes[i])
	}

	fmt.Println("Notes Seeded:")
}
