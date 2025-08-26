package app_db

import (
	"app/app_models"
	"fmt"
	"os"

	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) {

	// Check if number of users are > 0
	var count int64
	db.Model(&app_models.Users{}).Count(&count)
	if count > 0 {
		fmt.Println("Users already seeded")
		return
	}

	users := []app_models.Users{
		{Email: "dummy@dum.my",
			Password: "$2a$10$gafMnfiqCrocD54Pk8IqB.RwW4LikoTwzhyHV73fyP,.07KmzGTlW", // Password: appadmin
			IsAuth:   false,
			Note:     "Dummy user 1 for foreign key",
			Role:     "dummy"},
		
		{Email: "admin@app.loc",
			Password: "$2a$10$gafMnfiqCrocF54Pk8IqB.RwW3LikotwzhyHV40fyP..07KmyGTlW", // Password: appadmin
			IsAuth:   true,
			Note:     "Default Admin User",
			Role:     "admin"},

		{Email: "super@app.loc",
			Password: "$2a$10$7PcHWQCFWYX8K.k8GwI8WeSCc8s/Xyly.HM0FfyUm2BE8tGTcfOYy", // Psddword: superuser
			IsAuth:   false,
			Note:     "Default Superuser",
			Role:     "superuser"},

		{Email: "user@app.loc",
			Password: "$2a$10$UT.RGY8jJ7uhdsDaFzn0e.e5DXeXDgcITJf3olYpqsYATa/dYw1bK", // Password: password
			IsAuth:   false,
			Note:     "Default User",
			Role:     "user"},
	}

	for i := range users {
		db.Create(&users[i])
	}

	if len(users) != 4 {
		fmt.Println("Error seeding users")
		os.Exit(1)
	}

	fmt.Println("Users Seeded:", users[0].Email, users[1].Email, users[2].Email)

}

func seedMenus(db *gorm.DB) {

	// Check if the database is empty
	var count int64
	db.Model(&app_models.Menu{}).Count(&count)
	if count > 0 {
		fmt.Println("MnuDB already seeded")
		return
	}

	// Seed the database with initial data
	menus := []app_models.Menu{
		{Title: "Title0", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title1", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title2", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title3", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title4", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title5", Type: "sub", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title6", Type: "mnu", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title7", Type: "mnu", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title8", Type: "mnu", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
		{Title: "Title9", Type: "mnu", Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}}},
	}

	for i := range menus {
		db.Create(&menus[i])
	}

	fmt.Println("Menus and Menusub Seeded")

}

func seedStatus(db *gorm.DB) {
	// Check if number of status are > 0
	var count int64
	db.Model(&app_models.Status{}).Count(&count)
	if count > 0 {
		fmt.Println("Status already seeded")
		return
	}

	status := []app_models.Status{
		{Title: "Project", Stasub: []app_models.Stasub{
			{Name: "New", Type: "[D]"},
			{Name: "Active", Type: ""},
			{Name: "Paused", Type: ""},
			{Name: "Archived", Type: ""},
		}},
		{Title: "All", Stasub: []app_models.Stasub{ // Address, SIMCard, Email, Phone, Computer, Vehicle
			{Name: "New", Type: "[D]"},
			{Name: "Active", Type: ""},
			{Name: "Inactive", Type: ""},
			{Name: "Archived", Type: ""},
			{Name: "Deleted", Type: ""},
		}},
		{Title: "Status 4", Stasub: []app_models.Stasub{
			{Name: "New", Type: "[D]"},
			{Name: "Active", Type: ""},
			{Name: "Inactive", Type: ""},
		}},
	}

	for i := range status {
		db.Create(&status[i])
	}

	fmt.Println("Statuses Seeded")

}

func seedTypes(db *gorm.DB) {
	// Check if number of types are > 0
	var count int64
	db.Model(&app_models.Type{}).Count(&count)
	if count > 0 {
		fmt.Println("Types already seeded")
		return
	}

	types := []app_models.Type{
		{Title: "Project", Typsub: []app_models.Typsub{
			{Name: "Cover", Type: ""},
			{Name: "DLB", Type: ""},
			{Name: "Other", Type: "[D]"},
		}},
		{Title: "Person", Typsub: []app_models.Typsub{
			{Name: "New", Type: "[D]"},
			{Name: "Agent", Type: ""},
			{Name: "Contact", Type: ""},
			{Name: "Cover", Type: ""},
			{Name: "Other", Type: ""},
		}},
		{Title: "All", Typsub: []app_models.Typsub{ // Address, Email, Phone
			{Name: "Private", Type: ""},
			{Name: "work", Type: ""},
			{Name: "Other", Type: "[D]"},
		}},
		{Title: "Note", Typsub: []app_models.Typsub{
			{Name: "Contact", Type: ""},
			{Name: "Security", Type: ""},
			{Name: "Other", Type: "[D]"},
		}},
	}

	for i := range types {
		db.Create(&types[i])
	}
	fmt.Println("Types Seeded")

}
