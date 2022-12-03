package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/packages/token"
	"gorm.io/gorm"

	userRepository "github.com/snykk/golib_backend/datasources/databases/users"
	userUsecase "github.com/snykk/golib_backend/domains/users"
	userController "github.com/snykk/golib_backend/http/controllers/users"
	"github.com/snykk/golib_backend/packages/cache"
)

type usersRoutes struct {
	controller     userController.UserController
	router         *gin.Engine
	db             *gorm.DB
	authMiddleware gin.HandlerFunc
}

func NewUsersRoute(db *gorm.DB, jwtService token.JWTService, redisCache cache.RedisCache, ristrettoCache cache.RistrettoCache, router *gin.Engine, authMiddleware gin.HandlerFunc) *usersRoutes {
	// user route
	userRepository := userRepository.NewPostgreUserRepository(db)
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
		userRoute.POST("/change-password", r.controller.ChangePassword)
		userRoute.POST("/change-email", r.controller.ChangeEmail)
	}
}
