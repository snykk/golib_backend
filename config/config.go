package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        int
	Environment string
	Debug       bool

	DBHost     string
	DBPort     int
	DBDatabase string
	DBUsername string
	DBPassword string
	DBDsn      string

	JWTSecret  string
	JWTExpired int
	JWTIssuer  string

	OTPEmail    string
	OTPPassword string

	REDISHost     string
	REDISPassword string
	REDISExpired  int
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	// assign value
	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.DBHost = viper.GetString("DB_HOST")
	AppConfig.DBPort = viper.GetInt("DB_PORT")
	AppConfig.DBDatabase = viper.GetString("DB_DATABASE")
	AppConfig.DBUsername = viper.GetString("DB_USERNAME")
	AppConfig.DBPassword = viper.GetString("DB_PASSWORD")
	AppConfig.DBDsn = viper.GetString("DB_DSN")

	AppConfig.JWTSecret = viper.GetString("JWT_SECRET")
	AppConfig.JWTExpired = viper.GetInt("JWT_EXPIRED")
	AppConfig.JWTIssuer = viper.GetString("JWT_ISSUER")

	AppConfig.OTPEmail = viper.GetString("OTP_EMAIL")
	AppConfig.OTPPassword = viper.GetString("OTP_PASSWORD")

	AppConfig.REDISHost = viper.GetString("REDIS_HOST")
	AppConfig.REDISPassword = viper.GetString("REDIS_PASS")
	AppConfig.REDISExpired = viper.GetInt("REDIS_EXPIRED")

	// check
	if AppConfig.Port == 0 || AppConfig.Environment == "" || AppConfig.JWTSecret == "" || AppConfig.JWTExpired == 0 || AppConfig.JWTIssuer == "" || AppConfig.OTPEmail == "" || AppConfig.OTPPassword == "" || AppConfig.REDISHost == "" || AppConfig.REDISPassword == "" || AppConfig.REDISExpired == 0 {
		return errors.New("required variabel environment is empty")
	}

	switch AppConfig.Environment {
	case "development":
		if AppConfig.DBHost == "" || AppConfig.DBPort == 0 || AppConfig.DBDatabase == "" || AppConfig.DBUsername == "" || AppConfig.DBPassword == "" {
			return errors.New("required variabel environment is empty")
		}
	case "production":
		if AppConfig.DBDsn == "" {
			return errors.New("required variabel environment is empty")
		}
	}

	log.Println("[INIT] configuration loaded")
	return nil
}
