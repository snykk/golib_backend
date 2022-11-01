package config

import (
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

	JWTSecret  string
	JWTExpired int
	JWTIssuer  string

	OTPUsername string
	OTPPassword string
	OTPExpired  int
}

func InitializeAppConfig() {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.DBHost = viper.GetString("DB_HOST")
	AppConfig.DBPort = viper.GetInt("DB_PORT")
	AppConfig.DBDatabase = viper.GetString("DB_DATABASE")
	AppConfig.DBUsername = viper.GetString("DB_USERNAME")
	AppConfig.DBPassword = viper.GetString("DB_PASSWORD")

	AppConfig.JWTSecret = viper.GetString("JWT_SECRET")
	AppConfig.JWTExpired = viper.GetInt("JWT_EXPIRED")
	AppConfig.JWTIssuer = viper.GetString("JWT_ISSUER")

	AppConfig.OTPUsername = viper.GetString("OTP_USERNAME")
	AppConfig.OTPPassword = viper.GetString("OTP_PASSWORD")
	AppConfig.OTPExpired = viper.GetInt("OTP_EXPIRED")
	log.Println("[INIT] configuration loaded")
}
