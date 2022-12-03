package drivers

import (
	"errors"
	"fmt"
	"log"

	configEnv "github.com/snykk/golib_backend/config"
	bookRepository "github.com/snykk/golib_backend/datasources/databases/books"
	userRepository "github.com/snykk/golib_backend/datasources/databases/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConfigPostgreSQL struct {
	DB_Username string
	DB_Password string
	DB_Host     string
	DB_Port     int
	DB_Database string
	DB_DSN      string
}

func DbMigrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&bookRepository.Book{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&userRepository.User{})
	return
}

func (config *ConfigPostgreSQL) InitializeDatabasePostgreSQL() (*gorm.DB, error) {
	var dsn string

	if configEnv.AppConfig.Environment == "development" {
		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.DB_Host, config.DB_Port, config.DB_Database,
			config.DB_Username, config.DB_Password)
	} else if configEnv.AppConfig.Environment == "production" {
		dsn = config.DB_DSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("[INIT] failed connecting to PostgreSQL")
	}
	log.Println("[INIT] connected to PostgreSQL")

	err = DbMigrate(db)
	if err != nil {
		return nil, errors.New("[INIT] failed when running migration")
	}

	log.Println("[INIT] migration success")

	return db, nil
}
