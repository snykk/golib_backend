package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/utils/token"
	"gorm.io/gorm"

	"github.com/snykk/golib_backend/datasources/cache"
	userRepository "github.com/snykk/golib_backend/datasources/postgre/users"
	userUsecase "github.com/snykk/golib_backend/domains/users"
	userController "github.com/snykk/golib_backend/http/controllers/users"
)

type usersRoutes struct {
	controller     userController.UserController
	router         *gin.Engine
	db             *gorm.DB
	authMiddleware gin.HandlerFunc
}

func NewUsersRoute(db *gorm.DB, jwtService token.JWTService, redisCache cache.RedisCache, ristrettoCache cache.RistrettoCache, router *gin.Engine, authMiddleware gin.HandlerFunc) *usersRoutes {
	// user route
	userRepository := userRepository.NewUserRepository(db)
	userUsecase := userUsecase.NewUserUsecase(userRepository, jwtService)
	userController := userController.NewUserController(userUsecase, redisCache, ristrettoCache)
	return &usersRoutes{controller: userController, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *usersRoutes) UsersRoute() {
	// Auth
	authRoute := r.router.Group("auth")
	authRoute.POST("/login", r.controller.Login)
	authRoute.POST("/regis", r.controller.Regis)
	authRoute.POST("/send-otp", r.controller.SendOTP)
	authRoute.POST("/verif-otp", r.controller.VerifOTP)

	// Users
	userRoute := r.router.Group("users")
	userRoute.Use(r.authMiddleware)
	{
		userRoute.GET("", r.controller.GetAll)
		userRoute.GET("/:id", r.controller.GetById)
		userRoute.GET("/me", r.controller.GetUserData)
		userRoute.PUT("", r.controller.Update)
		userRoute.DELETE("", r.controller.Delete)
	}
}
