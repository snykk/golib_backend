package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/http/token"
	"gorm.io/gorm"

	"github.com/snykk/golib_backend/datasources/cache"
	reviewRepository "github.com/snykk/golib_backend/datasources/databases/reviews"
	reviewUseCase "github.com/snykk/golib_backend/domains/reviews"
	reviewController "github.com/snykk/golib_backend/http/controllers/reviews"
)

type reviewsRoutes struct {
	controller     reviewController.ReviewController
	router         *gin.Engine
	db             *gorm.DB
	authMiddleware gin.HandlerFunc
}

func NewReviewsRoute(db *gorm.DB, jwtService token.JWTService, ristrettoCache cache.RistrettoCache, router *gin.Engine, authMiddleware gin.HandlerFunc) *reviewsRoutes {
	reviewRepository := reviewRepository.NewPostgreReviewRepository(db)
	reviewUseCase := reviewUseCase.NewReviewUsecase(reviewRepository)
	reviewController := reviewController.NewReviewController(reviewUseCase, ristrettoCache)

	return &reviewsRoutes{controller: reviewController, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *reviewsRoutes) ReviewsRoute() {
	// => Review
	reviewRoute := r.router.Group("reviews")
	reviewRoute.Use(r.authMiddleware)
	{
		reviewRoute.POST("", r.controller.Store)
		reviewRoute.GET("", r.controller.GetAll)
		reviewRoute.GET("/:id", r.controller.GetById)
		reviewRoute.GET("/book/:id", r.controller.GetByBookId)
		reviewRoute.GET("/user/:id", r.controller.GetByUserid)
		reviewRoute.PUT("/:id", r.controller.Update)
		reviewRoute.DELETE("/:id", r.controller.Delete)
	}

}
