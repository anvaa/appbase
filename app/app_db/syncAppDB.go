package app_db

import (
	"gorm.io/gorm"

	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"app/app_models"
)

func SyncAppDB(db *gorm.DB) {
	db.AutoMigrate(
		&app_models.Users{},
		&app_models.AuthLevel{},

		&app_models.Status{},
		&app_models.Stasub{},
		&app_models.Menu{},
		&app_models.Menusub{},
		&app_models.Type{},
		&app_models.Typsub{},
	)

	seedMenus(db)
	seedStatus(db)
	seedTypes(db)

	seedAuthLevels(db)
	seedUsers(db)

	fmt.Println("Database Migrated")
}

func seedAuthLevels(db *gorm.DB) {
	var count int64
	db.Model(&app_models.AuthLevel{}).Count(&count)
	if count > 0 {
		fmt.Println("AuthLevels already seeded")
		return
	}

	authLevels := []app_models.AuthLevel{
		{Name: "admin", Level: 10},
		{Name: "super", Level: 5},
		{Name: "user", Level: 1},
	}

	for _, a := range authLevels {
		db.Create(&a)
	}

	fmt.Println("AuthLevels Seeded")
}

func seedUsers(db *gorm.DB) {
	// Check if number of users are > 0
	var count int64
	db.Model(&app_models.Users{}).Count(&count)
	if count > 0 {
		fmt.Println("Users already seeded")
		return
	}

	// Plaintext passwords for seeding
	userSeeds := []struct {
		Email       string
		Password    string
		IsAuth      bool
		AuthLevelID int
		Note        string
		Orgname     string
	}{
		{"admin@app.loc", "appadmin", true, 1, "Default Admin User", "app"},
		{"super@app.loc", "superuser", false, 2, "Default Superuser", "app"},
		{"user@app.loc", "password", false, 3, "Default User", "app"},
	}

	var users []app_models.Users
	for _, u := range userSeeds {
		hashed, err := hashPassword(u.Password)
		if err != nil {
			fmt.Println("Error hashing password for", u.Email, ":", err)
			os.Exit(1)
		}
		users = append(users, app_models.Users{
			Email:       u.Email,
			Password:    hashed,
			IsAuth:      u.IsAuth,
			Note:        u.Note,
			AuthLevelID: u.AuthLevelID,
			Orgname:     u.Orgname,
		})
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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func seedMenus(db *gorm.DB) {
	var count int64
	db.Model(&app_models.Menu{}).Count(&count)
	if count > 0 {
		fmt.Println("Menus already seeded")
		return
	}

	menuDefs := []struct {
		Title string
		Type  string
	}{
		{"Title1", "sub"},
		{"Title2", "sub"},
		{"Title3", "sub"},
		{"Title4", "sub"},
		{"Title5", "sub"},
		{"Title6", "mnu"},
	}

	menus := make([]app_models.Menu, len(menuDefs))
	for i, def := range menuDefs {
		menus[i] = app_models.Menu{
			Title:   def.Title,
			Type:    def.Type,
			Menusub: []app_models.Menusub{{Name: "def", Type: "[D]"}},
		}
	}

	for i := range menus {
		db.Create(&menus[i])
	}

	fmt.Println("Menus and Menusub Seeded")
}

func seedStatus(db *gorm.DB) {
	var count int64
	db.Model(&app_models.Status{}).Count(&count)
	if count > 0 {
		fmt.Println("Status already seeded")
		return
	}

	statuses := []struct {
		Title  string
		Stasub []app_models.Stasub
	}{
		{"Status 1", []app_models.Stasub{
			{Name: "New", Type: "[D]"},
			{Name: "Other 11", Type: ""},
		}},
		{"Status 2", []app_models.Stasub{
			{Name: "New", Type: "[D]"},
			{Name: "Other 22", Type: ""},
		}},
		{"Status 3", []app_models.Stasub{
			{Name: "New", Type: "[D]"},
			{Name: "Other 33", Type: ""},
		}},
	}

	for _, s := range statuses {
		status := app_models.Status{
			Title:  s.Title,
			Stasub: s.Stasub,
		}
		db.Create(&status)
	}

	fmt.Println("Statuses Seeded")
}

func seedTypes(db *gorm.DB) {
	var count int64
	db.Model(&app_models.Type{}).Count(&count)
	if count > 0 {
		fmt.Println("Types already seeded")
		return
	}

	typeDefs := []struct {
		Title  string
		Typsub []app_models.Typsub
	}{
		{"Type 1", []app_models.Typsub{
			{Name: "type11", Type: "[D]"},
			{Name: "type12", Type: ""},
		}},
		{"Type 2", []app_models.Typsub{
			{Name: "type21", Type: "[D]"},
			{Name: "type22", Type: ""},
		}},
		{"Type 3", []app_models.Typsub{
			{Name: "type31", Type: "[D]"},
			{Name: "type32", Type: ""},
		}},
	}

	for _, def := range typeDefs {
		t := app_models.Type{
			Title:  def.Title,
			Typsub: def.Typsub,
		}
		db.Create(&t)
	}
	fmt.Println("Types Seeded")
}
