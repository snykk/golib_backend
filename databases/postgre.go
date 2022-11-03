package postgres

import (
	"fmt"
	"log"

	configEnv "github.com/snykk/golib_backend/config"
	bookRepository "github.com/snykk/golib_backend/databases/books"
	userRepository "github.com/snykk/golib_backend/databases/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConfigDB struct {
	DB_Username string
	DB_Password string
	DB_Host     string
	DB_Port     int
	DB_Database string
}

func DbMigrate(db *gorm.DB) {
	db.AutoMigrate(&bookRepository.Book{})
	db.AutoMigrate(&userRepository.User{})
}

func (config *ConfigDB) InitializeDatabase() *gorm.DB {
	var dsn string

	if configEnv.AppConfig.Environment == "development" {
		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.DB_Host, config.DB_Port, config.DB_Database,
			config.DB_Username, config.DB_Password)
	} else if configEnv.AppConfig.Environment == "production" {
		dsn = configEnv.AppConfig.DBDsn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("[INIT] failed connecting to PostgreSQL")
		return nil
	}
	log.Println("[INIT] connected to PostgreSQL")

	DbMigrate(db)

	log.Println("[INIT] migration success")

	return db
}
