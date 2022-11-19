package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/http/middlewares"
	"github.com/snykk/golib_backend/utils/token"
	"gorm.io/gorm"

	bookRepository "github.com/snykk/golib_backend/datasources/postgre/books"
	userRepository "github.com/snykk/golib_backend/datasources/postgre/users"
	bookUsecase "github.com/snykk/golib_backend/domains/books"
	userUsecase "github.com/snykk/golib_backend/domains/users"
	bookController "github.com/snykk/golib_backend/http/controllers/books"
	userController "github.com/snykk/golib_backend/http/controllers/users"
)

func InitializeRouter(conn *gorm.DB) (router *gin.Engine) {
	router = gin.Default()

	// middleware jwt
	jwtService := token.NewJWTService()

	// CACHE
	redisCache := cache.NewRedisCache(config.AppConfig.REDISHost, 0, config.AppConfig.REDISPassword, time.Duration(config.AppConfig.REDISExpired))
	ristrettoCache, err := cache.NewRistrettoCache()
	if err != nil {
		panic(err)
	}

	// user route
	userRepository := userRepository.NewUserRepository(conn)
	userUsecase := userUsecase.NewUserUsecase(userRepository, jwtService)
	userController := userController.NewUserController(userUsecase, redisCache, ristrettoCache)

	// book route
	bookRepository := bookRepository.NewBookRepository(conn)
	bookUsecase := bookUsecase.NewBookUsecase(bookRepository)
	bookController := bookController.NewBookController(bookUsecase, ristrettoCache)

	// ===============> LIST OF ROUTE <==================
	// => Root
	router.GET("/", rootHandler)

	// => Auth
	authRoute := router.Group("auth")
	authRoute.POST("/login", userController.Login)
	authRoute.POST("/regis", userController.Regis)
	authRoute.POST("/send-otp", userController.SendOTP)
	authRoute.POST("/verif-otp", userController.VerifOTP)

	// => User
	userRoute := router.Group("users")
	userRoute.Use(middlewares.AuthorizeJWT(jwtService))
	{
		userRoute.GET("", userController.GetAll)
		userRoute.GET("/:id", userController.GetById)

		// encapsulate action for each user
		userRoute.Use(middlewares.IsValidUser(jwtService))
		{
			userRoute.PUT("/:id", userController.Update)
			userRoute.DELETE("/:id", userController.Delete)
		}
	}

	// => Book
	bookRoute := router.Group("books")
	bookRoute.Use(middlewares.AuthorizeJWT(jwtService))
	{
		bookRoute.GET("", bookController.GetAll)
		bookRoute.GET("/:id", bookController.GetById)

		// admin middleware
		bookRoute.Use(middlewares.IsAdmin(jwtService))
		{
			bookRoute.POST("", bookController.Store)
			bookRoute.PUT("/:id", bookController.Update)
			bookRoute.DELETE("/:id", bookController.Delete)
		}
	}

	log.Println("[INIT] router success")
	return
}

type Base struct {
	Routes     Routes            `json:"routes"`
	Middleware map[string]string `json:"middleware"`
	Maintainer string            `json:"maintainer"`
	Repository string            `json:"repository"`
}

type Routes struct {
	Auth  map[string]string `json:"auth"`
	Users map[string]string `json:"users"`
	Books map[string]string `json:"books"`
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Base{
		Routes: Routes{
			Auth: map[string]string{
				"Login [POST]":     "/auth/login",
				"Regis [POST]":     "/auth/regis",
				"Send OTP [POST]":  "/auth/send-otp",
				"Verif OTP [POST]": "/auth/verif-otp",
			},
			Users: map[string]string{
				"Get Users [GET] <AuthorizeJWT>":                 "/users",
				"Get User [GET] <AuthorizeJWT>":                  "/users/:id",
				"Update User [PUT] <AuthorizeJWT> <IsValidUser>": "/users/:id",
				"Delete User [PUT] <AuthorizeJWT> <IsValidUser>": "/users/:id",
			},
			Books: map[string]string{
				"Get Books [GET] <AuthorizeJWT>":              "/books",
				"Get Book [GET] <AuthorizeJWT>":               "/books/:id",
				"Create Book [POST] <AuthorizeJWT> <IsAdmin>": "/books",
				"Update Book [PUT] <AuthorizeJWT> <IsAdmin>":  "/books/:id",
				"Delete Book [PUT] <AuthorizeJWT> <IsAdmin>":  "/books/:id",
			},
		},
		Middleware: map[string]string{
			"<AuthorizeJWT>": "only user with valid token can access endpoint",
			"<IsValidUser>":  "only user itself or admin can access endpoint",
			"<IsAdmin>":      "only admin can access endpoint",
		},
		Maintainer: "Moh. Najib Fikri aka snykk github.com/snykk najibfikri13@gmail.com",
		Repository: "https://github.com/snykk/golib-backend",
	})
}
