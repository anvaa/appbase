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
		{Name: "Project Alpha", Note: "First project", StasubID : 1, TypsubID: 1, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
		{Name: "Project Beta", Note: "", StasubID: 2, TypsubID: 2, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
	}

	for i := range projects {
		db.Create(&projects[i])
	}

	fmt.Println("Projects Seeded:")
}

func seedEmails(db *gorm.DB ,email string, projid int) app_models.Email {
	emailModel := app_models.Email{
		
		Address:   email,
		ProjectID: projid,
	}

	if err := db.Create(&emailModel).Error; err != nil {
		fmt.Println("Error seeding email:", err)
		return app_models.Email{}
	}

	fmt.Println("Email Seeded:", emailModel.Address)
	return emailModel
}	