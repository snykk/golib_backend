package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/utils/token"
	"gorm.io/gorm"

	"github.com/snykk/golib_backend/datasources/cache"
	bookRepository "github.com/snykk/golib_backend/datasources/postgre/books"
	bookUseCase "github.com/snykk/golib_backend/domains/books"
	bookController "github.com/snykk/golib_backend/http/controllers/books"
)

type booksRoutes struct {
	controller          bookController.BookController
	router              *gin.Engine
	db                  *gorm.DB
	authMiddleware      gin.HandlerFunc
	authAdminMiddleware gin.HandlerFunc
}

func NewBooksRoute(db *gorm.DB, jwtService token.JWTService, ristrettoCache cache.RistrettoCache, router *gin.Engine, authMiddleware gin.HandlerFunc, authAdminMiddleware gin.HandlerFunc) *booksRoutes {
	// user route
	bookRepository := bookRepository.NewBookRepository(db)
	bookUseCase := bookUseCase.NewBookUsecase(bookRepository)
	bookController := bookController.NewBookController(bookUseCase, ristrettoCache)
	return &booksRoutes{controller: bookController, router: router, db: db, authMiddleware: authMiddleware, authAdminMiddleware: authAdminMiddleware}
}

func (r *booksRoutes) BooksRoute() {
	// => Book

	// all users
	bookRoute := r.router.Group("books")
	bookRoute.Use(r.authMiddleware)
	{
		bookRoute.GET("", r.controller.GetAll)
		bookRoute.GET("/:id", r.controller.GetById)
	}

	// admin endpoint
	bookAdminRoute := r.router.Group("books")
	bookAdminRoute.Use(r.authAdminMiddleware)
	{
		bookRoute.POST("", r.controller.Store)
		bookRoute.PUT("/:id", r.controller.Update)
		bookRoute.DELETE("/:id", r.controller.Delete)
	}
}
