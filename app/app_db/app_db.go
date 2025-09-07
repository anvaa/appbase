package app_db

import (
	"log"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"app/app_models"
)

var AppDB *gorm.DB

func CnnAppDB(config app_models.DbConfig) {

	switch config.Type {
		case "sqlite":
			cnnSqlite(config.Path)
		case "mysql":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				config.User, config.Password, config.Host, config.Port, config.Name)
			cnnMysql(dsn)
		case "postgres":
			dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
				config.Host, config.Port, config.User, config.Name, config.Password)
			cnnPostgres(dsn)
		default:
			log.Fatalf("unsupported database type: %s", config.Type)
	}	

	SyncAppDB(AppDB)
	log.Println(config.Type, "connected and synced")
}

func CloseAppDB() {
	sqlDB, err := AppDB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("failed to close database: %v", err)
	}
	log.Println("AppDB closed")
}

func cnnSqlite(dbpath string) {
	var err error
	AppDB, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", dbpath, err)
	}

	if err := AppDB.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		log.Fatalf("failed to set foreign_keys pragma: %v", err)
	}

	// Uncomment the following line if you want to enable WAL mode
	// if err := AppDB.Exec("PRAGMA journal_mode=WAL;").Error; err != nil {
	// 	log.Fatalf("failed to set journal_mode pragma: %v", err)
	// }

}

func cnnMysql(dsn string) {
	var err error
	AppDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to mysql: %v", err)
	}
	
}

func cnnPostgres(dsn string) {
	var err error
	AppDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	
}
