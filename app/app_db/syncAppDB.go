package app_db

import (
	"gorm.io/gorm"

	"fmt"
	"os"

	"app/app_models"
)

func SyncAppDB(db *gorm.DB) {
	db.AutoMigrate(
		&app_models.Users{},

		&app_models.Status{},
		&app_models.StaSub{},
		&app_models.Menu{},
		&app_models.MenuSub{},
		// &app_models.StatusHistory{},
	)

	seedUsers(db)
	seedMenus(db)
	seedStatus(db)
}

func seedUsers(db *gorm.DB) {

	// Check if number of users are > 0
	var count int64
	db.Model(&app_models.Users{}).Count(&count)
	if count > 0 {
		fmt.Println("Users already seeded")
		return
	}

	users := []app_models.Users{
		{Email: "admin@app.loc",
			Password: "$2a$10$gafMnfiqCrocF54Pk8IqB.RwW3LikotwzhyHV40fyP..07KmyGTlW", // Password: appadmin
			IsAuth:   true,
			Note:     "Default Admin User",
			Role:     "admin"},

		{Email: "super@app.loc",
			Password: "$2a$10$7PcHWQCFWYX8K.k8GwI8WeSCc8s/Xyly.HM0FfyUm2BE8tGTcfOYy", // Psddword: superuser
			IsAuth:   false,
			Note:     "Default SuperUser",
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

	if len(users) != 3 {
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
		{Title: "Title0", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title1", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title2", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title3", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title4", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title5", Type: "sub", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title6", Type: "mnu", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title7", Type: "mnu", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title8", Type: "mnu", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
		{Title: "Title9", Type: "mnu", MenuSub: []app_models.MenuSub{{Name: "def", Type: "[D]"}}},
	}

	for i := range menus {
		db.Create(&menus[i])
	}

	fmt.Println("Menus and MenuSub Seeded")

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
		{Title: "Status0", StaSub: []app_models.StaSub{
			{Name: "Ny", Type: "[D]"},
			{Name: "Notat0", Type: ""},
		}},
		{Title: "Status1", StaSub: []app_models.StaSub{
			{Name: "Ny", Type: "[D]"},
			{Name: "Notat1", Type: ""},
		}},
		{Title: "Status2", StaSub: []app_models.StaSub{
			{Name: "Ny", Type: "[D]"},
			{Name: "Notat2", Type: ""},
		}},
	}

	for i := range status {
		db.Create(&status[i])
	}

	fmt.Println("Statuses Seeded")

}
