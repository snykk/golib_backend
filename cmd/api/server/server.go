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

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/datasources/databases/drivers"
	"github.com/snykk/golib_backend/http/logger"
	"github.com/snykk/golib_backend/http/middlewares"
	"github.com/snykk/golib_backend/http/routes"
	"github.com/snykk/golib_backend/http/token"
	"gorm.io/gorm"
)

type App struct {
	httpServer *http.Server
}

func NewApp() (*App, error) {
	// setup databases
	conn, err := setupDatabse()
	if err != nil {
		return nil, err
	}

	// setup router
	router := setupRouter()

	// jwt service
	jwtService := token.NewJWTService()

	// cache
	redisCache := cache.NewRedisCache(config.AppConfig.REDISHost, 0, config.AppConfig.REDISPassword, time.Duration(config.AppConfig.REDISExpired))
	ristrettoCache, err := cache.NewRistrettoCache()
	if err != nil {
		panic(err)
	}

	// user middleware
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, false)
	// admin middleware
	authAdminMiddleware := middlewares.NewAuthMiddleware(jwtService, true)

	// Routes
	router.GET("/", routes.RootHandler)
	routes.NewUsersRoute(conn, jwtService, redisCache, ristrettoCache, router, authMiddleware).UsersRoute()
	routes.NewBooksRoute(conn, jwtService, ristrettoCache, router, authMiddleware, authAdminMiddleware).BooksRoute()
	routes.NewReviewsRoute(conn, jwtService, ristrettoCache, router, authMiddleware).ReviewsRoute()

	// setup http server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		httpServer: server,
	}, nil
}

func (a *App) Run() error {
	// Gracefull Shutdown
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

func setupDatabse() (*gorm.DB, error) {
	// Setup Config Databse
	configDB := drivers.ConfigPostgreSQL{
		DB_Username: config.AppConfig.DBUsername,
		DB_Password: config.AppConfig.DBPassword,
		DB_Host:     config.AppConfig.DBHost,
		DB_Port:     config.AppConfig.DBPort,
		DB_Database: config.AppConfig.DBDatabase,
		DB_DSN:      config.AppConfig.DBDsn,
	}

	// Initialize Database driversSQL
	conn, err := configDB.InitializeDatabasePostgreSQL()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func setupRouter() *gin.Engine {
	// set the runtime mode
	var mode = gin.ReleaseMode
	if config.AppConfig.Debug {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// create a new router instance
	router := gin.New()

	// set up middlewares
	router.Use(middlewares.CORSMiddleware())
	if mode == gin.DebugMode {
		router.Use(gin.LoggerWithFormatter(logger.CustomLogFormatter))
	}
	router.Use(gin.Recovery())

	return router
}
