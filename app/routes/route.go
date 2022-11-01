package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/app/middlewares"
	"github.com/snykk/golib_backend/cache"
	"github.com/snykk/golib_backend/utils/token"
	"gorm.io/gorm"

	bookController "github.com/snykk/golib_backend/controllers/books"
	userController "github.com/snykk/golib_backend/controllers/users"
	bookRepository "github.com/snykk/golib_backend/databases/books"
	userRepository "github.com/snykk/golib_backend/databases/users"
	bookUsecase "github.com/snykk/golib_backend/usecase/books"
	userUsecase "github.com/snykk/golib_backend/usecase/users"
)

func InitializeRouter(conn *gorm.DB) (router *gin.Engine) {
	router = gin.Default()

	// middleware jwt
	jwtService := token.NewJWTService()

	// CACHE
	redisCache := cache.NewRedisCache("localhost:6379", 0, 5)
	// user route
	userRepository := userRepository.NewUserRepository(conn)
	userUsecase := userUsecase.NewUserUsecase(userRepository, jwtService)
	userController := userController.NewUserController(userUsecase, redisCache)

	// book route
	bookRepository := bookRepository.NewBookRepository(conn)
	bookUsecase := bookUsecase.NewBookUsecase(bookRepository)
	bookController := bookController.NewBookController(bookUsecase)

	// ===============> LIST OF ROUTE <==================
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
