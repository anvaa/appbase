package app_db

import (
	"app/app_models"

	"fmt"
	"time"

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
		{Name: "Dummy Project", StasubID: 1, TypsubID: 1},
		{Name: "Project Alpha", StasubID: 1, TypsubID: 1},
		{Name: "Project Beta", StasubID: 2, TypsubID: 2},
	}

	for i := range projects {
		db.Create(&projects[i])
	}

	fmt.Println("Projects Seeded:")
}

func seedPerson(db *gorm.DB) {
	// Check if the number of persons is > 0
	var count int64
	db.Model(&app_models.Person{}).Count(&count)
	if count > 0 {
		fmt.Println("Persons already seeded")
		return
	}

	persons := []app_models.Person{
		{FirstName: "Joe", LastName: "Doe", DOB: time.Now(), Gender: "Male", Nationality: "Norwegian", ProjectID: 1, StasubID: 1, TypsubID: 1},
		{FirstName: "John", LastName: "Doe", DOB: time.Now(), Gender: "Male", Nationality: "ACanadian", ProjectID: 2, StasubID: 1, TypsubID: 1},
		{FirstName: "Jane", LastName: "Smith", DOB: time.Now(), Gender: "Female", Nationality: "British", ProjectID: 3, StasubID: 2, TypsubID: 2},
	}

	for i := range persons {
		db.Create(&persons[i])
	}

	fmt.Println("Persons Seeded:")
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
		{Header: "Dummy Note", Content: "Initial note for Dummy Project and Person",
			ProjectID: 1, PersonID: 1, TypsubID: 1},
		{Header: "Initial Note Alfa", Content: "Initial note for Project Alpha",
			ProjectID: 2, PersonID: 1, TypsubID: 7},
		{Header: "Second Note Alfa", Content: "Second note for Project Alpha",
			ProjectID: 2, PersonID: 2, TypsubID: 7},
		{Header: "Initial Note Beta", Content: "Initial note for Project Beta",
			ProjectID: 3, PersonID: 3, TypsubID: 7},
	}

	if err := db.Create(&notes).Error; err != nil {
		fmt.Printf("Failed to seed notes: %v\n", err)
		return
	}

	fmt.Println("Notes Seeded:")
}

func seedEmail(db *gorm.DB) {
	// Check if the number of emails is > 0
	var count int64
	db.Model(&app_models.Emails{}).Count(&count)
	if count > 0 {
		fmt.Println("Emails already seeded")
		return
	}

	emails := []app_models.Emails{
		{Email: "dummy@du.mmy", Password: "password123", Nationality: "Norwegian",
			PersonID: 1, ProjectID: 1, TypsubID: 11, StasubID: 5},
		{Email: "john.doe@example.com", Password: "password123", Nationality: "Canadian",
			PersonID: 2, ProjectID: 2, TypsubID: 11, StasubID: 5},
		{Email: "jane.smith@example.com", Password: "password123", Nationality: "British",
			PersonID: 3, ProjectID: 3, TypsubID: 11, StasubID: 5},
		{Email: "jenny.hansen@example.com", Password: "password123", Nationality: "British",
			PersonID: 1, ProjectID: 2, TypsubID: 11, StasubID: 5},
		{Email: "peter@example.com", Password: "password123", Nationality: "British",
			PersonID: 3, ProjectID: 3, TypsubID: 11, StasubID: 5},
	}

	for i := range emails {
		db.Create(&emails[i])
	}

	fmt.Println("Emails Seeded:")
}
