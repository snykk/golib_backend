package drivers

import (
	"errors"
	"fmt"
	"log"
	"time"

	configEnv "github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/constants"
	bookRepository "github.com/snykk/golib_backend/datasources/databases/books"
	userRepository "github.com/snykk/golib_backend/datasources/databases/users"
	"github.com/snykk/golib_backend/helpers"
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

func dbMigrate(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&bookRepository.Book{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&userRepository.User{})
	return
}

func (config *ConfigPostgreSQL) InitializeDatabasePostgreSQL() (*gorm.DB, error) {
	var dsn string

	if configEnv.AppConfig.Environment == constants.EnvironmentDevelopment {
		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.DB_Host, config.DB_Port, config.DB_Database,
			config.DB_Username, config.DB_Password)
	} else if configEnv.AppConfig.Environment == constants.EnvironmentProduction {
		dsn = config.DB_DSN
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("[INIT] failed connecting to PostgreSQL")
	}
	log.Println("[INIT] connected to PostgreSQL")

	if configEnv.AppConfig.Environment == constants.EnvironmentDevelopment {
		if err = db.Migrator().DropTable("users", "roles", "genders", "books"); err != nil {
			return nil, errors.New("[INIT ]failed droping tables:" + err.Error())
		}
		log.Println("[INIT] droping tables success")
	}

	err = dbMigrate(db)
	if err != nil {
		return nil, errors.New("[INIT] failed when running migration")
	}

	log.Println("[INIT] migration success")

	if configEnv.AppConfig.Environment == constants.EnvironmentDevelopment {
		lazySeeder(db)
		log.Println("[INIT] lazy seeders success")
	}

	return db, nil
}

func lazySeeder(db *gorm.DB) (err error) {
	// Role
	role1 := userRepository.Role{
		Id:   1,
		Name: "admin",
	}

	err = db.Model(&userRepository.Role{}).Create(&role1).Error
	if err != nil {
		return
	}

	role2 := userRepository.Role{
		Id:   2,
		Name: "user",
	}

	err = db.Model(&userRepository.Role{}).Create(&role2).Error
	if err != nil {
		return
	}

	// Gender
	gender1 := userRepository.Gender{
		Id:   1,
		Name: "male",
	}

	err = db.Model(&userRepository.Gender{}).Create(&gender1).Error
	if err != nil {
		return
	}

	gender2 := userRepository.Gender{
		Id:   2,
		Name: "female",
	}

	err = db.Model(&userRepository.Gender{}).Create(&gender2).Error
	if err != nil {
		return
	}

	// User
	userPass, _ := helpers.GenerateHash("12345")
	users1 := userRepository.User{
		Id:          1,
		FullName:    "patrict star",
		Username:    "pStar7",
		Email:       "najibfikri13@gmail.com",
		Password:    userPass,
		IsActivated: true,
		RoleId:      1,
		GenderId:    1,
		CreatedAt:   time.Now(),
	}
	err = db.Model(&userRepository.User{}).Create(&users1).Error
	if err != nil {
		return
	}

	return
}
