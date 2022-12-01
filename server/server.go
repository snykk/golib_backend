package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/snykk/golib_backend/config"
	postgre "github.com/snykk/golib_backend/datasources/postgre"
	"github.com/snykk/golib_backend/http/routes"
)

type App struct {
	httpServer *http.Server
}

func NewApp() (*App, error) {
	// Setup Config Databse
	configDB := postgre.ConfigDB{
		DB_Username: config.AppConfig.DBUsername,
		DB_Password: config.AppConfig.DBPassword,
		DB_Host:     config.AppConfig.DBHost,
		DB_Port:     config.AppConfig.DBPort,
		DB_Database: config.AppConfig.DBDatabase,
		DB_DSN:      config.AppConfig.DBDsn,
	}

	// Initialize Database PostgreSQL
	conn, err := configDB.InitializeDatabase()
	if err != nil {
		return nil, err
	}

	// setup http server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        routes.InitializeRouter(conn),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		httpServer: server,
	}, nil
}

func (a *App) Run() error {
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// make blocking channel and waiting for a signal
	<-quit
	log.Println("[CLOSE] shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("[CLOSE] error when shutdown server: %v", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("[CLOSE] timeout of 5 seconds.")
	log.Println("[CLOSE] server exiting")
	return nil
}
