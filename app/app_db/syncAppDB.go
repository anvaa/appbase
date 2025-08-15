package app_db

import (
	"gorm.io/gorm"

	"app/app_models"
)

func SyncAppDB(db *gorm.DB) {
	db.AutoMigrate(
		// Base models
		&app_models.Users{},

		&app_models.Status{},
		&app_models.Stasub{},
		&app_models.Menu{},
		&app_models.Menusub{},
		&app_models.Type{},
		&app_models.Typsub{},

		// App models
		&app_models.Project{},
		&app_models.Notes{},
	)

	// Seed the database with initial data
	// BaseApp models
	seedUsers(db)
	seedMenus(db)
	seedStatus(db)
	seedTypes(db)

	// App models
	seedProject(db)
	seedNotes(db)
}

