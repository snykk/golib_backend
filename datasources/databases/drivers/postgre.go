package drivers

import (
	"errors"
	"fmt"
	"log"
	"time"

	configEnv "github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/constants"
	bookRepository "github.com/snykk/golib_backend/datasources/databases/books"
	reviewRepository "github.com/snykk/golib_backend/datasources/databases/reviews"
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
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&reviewRepository.Review{})
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
		if err = db.Migrator().DropTable("users", "roles", "genders", "books", "reviews"); err != nil {
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
	// user1
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

	// user2
	users2 := userRepository.User{
		Id:          2,
		FullName:    "someone",
		Username:    "someonee",
		Email:       "someone@gmail.com",
		Password:    userPass,
		IsActivated: true,
		RoleId:      2,
		GenderId:    2,
		CreatedAt:   time.Now(),
	}
	err = db.Model(&userRepository.User{}).Create(&users2).Error
	if err != nil {
		return
	}

	// Book
	// book1
	book1 := bookRepository.Book{
		Id:          1,
		Title:       "Atomic Habits",
		Description: "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Voluptas cum quas veritatis voluptatem quia id voluptates, eum voluptatum officiis sed, maxime reprehenderit aut, magnam illo architecto earum consectetur ipsam a.",
		Author:      "James Clear",
		Publisher:   "Gramedia",
		ISBN:        "1234567891234",
		Rating:      9,
		CreatedAt:   time.Now(),
	}
	err = db.Model(&bookRepository.Book{}).Create(&book1).Error
	if err != nil {
		return
	}

	// book2
	book2 := bookRepository.Book{
		Id:          2,
		Title:       "Mindset",
		Description: "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Voluptas cum quas veritatis voluptatem quia id voluptates, eum voluptatum officiis sed, maxime reprehenderit aut, magnam illo architecto earum consectetur ipsam a.",
		Author:      "Carrol Dweck",
		Publisher:   "Gramedia",
		ISBN:        "1234567891234",
		CreatedAt:   time.Now(),
	}
	err = db.Model(&bookRepository.Book{}).Create(&book2).Error
	if err != nil {
		return
	}

	// Book

	// book1
	review1 := reviewRepository.Review{
		Id:        1,
		Text:      "ini review",
		Rating:    9,
		BookId:    1,
		UserId:    1,
		CreatedAt: time.Now(),
	}
	err = db.Model(&reviewRepository.Review{}).Create(&review1).Error
	if err != nil {
		return
	}

	return
}
