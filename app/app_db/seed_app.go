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
		{Name: "Project Alpha", StasubID: 1, TypsubID: 1, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
		{Name: "Project Beta", StasubID: 2, TypsubID: 2, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
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
		{FirstName: "John", LastName: "Doe", DOB: time.Now(), Gender: "Male", Nationality: "American", ProjectID: 1, StasubID: 1, TypsubID: 1, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
		{FirstName: "Jane", LastName: "Smith", DOB: time.Now(), Gender: "Female", Nationality: "British", ProjectID: 1, StasubID: 2, TypsubID: 2, CreatedbyID: 1, UpdatedbyID: 1, DeletedbyID: 1},
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
		{Header: "Note Header 1", Content: "Note 1 Note 1 Note 1 Note 1 Note 1 Note 1 Note 1 Note 1 Note 1 Note 1 ", ProjectID: 1, PersonID: 1, TypsubID: 7},
		{Header: "Note Header 2 Note Header 2", Content: "Note 2 Note 2 Note 2 Note 2 Note 2 Note 2 Note 2 Note 2 Note 2 Note 2 ", ProjectID: 1, PersonID: 1, TypsubID: 8},
		{Header: "Note Header 3 Note 3", Content: "Note 3 Note 3 Note 3 Note 3 Note 3 Note 3 Note 3 Note 3 Note 3 Note 3 ", ProjectID: 1, PersonID: 1, TypsubID: 9},
		{Header: "Note Header 4", Content: "Note 4 Note 4 Note 4 Note 4 Note 4 Note 4 Note 4 Note 4 Note 4 Note 4 ", ProjectID: 2, PersonID: 1, TypsubID: 7},
		{Header: "Note Header 5 Note Header 5", Content: "Note 5 Note 5 Note 5 Note 5 Note 5 Note 5 Note 5 Note 5 Note 5 Note 5 ", ProjectID: 2, PersonID: 1, TypsubID: 8},
	}

	for i := range notes {
		db.Create(&notes[i])
	}

	fmt.Println("Notes Seeded:")
}
