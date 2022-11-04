package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/app/routes"
	"github.com/snykk/golib_backend/config"

	postgre "github.com/snykk/golib_backend/datasources/postgre"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		log.Fatal(err)
	}

	if !config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	configDB := postgre.ConfigDB{
		DB_Username: config.AppConfig.DBUsername,
		DB_Password: config.AppConfig.DBPassword,
		DB_Host:     config.AppConfig.DBHost,
		DB_Port:     config.AppConfig.DBPort,
		DB_Database: config.AppConfig.DBDatabase,
		DB_DSN:      config.AppConfig.DBDsn,
	}

	conn := configDB.InitializeDatabase()

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        routes.InitializeRouter(conn),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
