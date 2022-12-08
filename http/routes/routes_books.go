package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/packages/token"
	"gorm.io/gorm"

	bookRepository "github.com/snykk/golib_backend/datasources/databases/books"
	bookUseCase "github.com/snykk/golib_backend/domains/books"
	bookController "github.com/snykk/golib_backend/http/controllers/books"
	"github.com/snykk/golib_backend/packages/cache"
)

type booksRoutes struct {
	controller          bookController.BookController
	router              *gin.Engine
	db                  *gorm.DB
	authMiddleware      gin.HandlerFunc
	authAdminMiddleware gin.HandlerFunc
}

func NewBooksRoute(db *gorm.DB, jwtService token.JWTService, ristrettoCache cache.RistrettoCache, router *gin.Engine, authMiddleware gin.HandlerFunc, authAdminMiddleware gin.HandlerFunc) *booksRoutes {
	bookRepository := bookRepository.NewPostgreBookRepository(db)
	bookUseCase := bookUseCase.NewBookUsecase(bookRepository)
	bookController := bookController.NewBookController(bookUseCase, ristrettoCache)

	return &booksRoutes{controller: bookController, router: router, db: db, authMiddleware: authMiddleware, authAdminMiddleware: authAdminMiddleware}
}

func (r *booksRoutes) BooksRoute() {
	// Book
	bookRoute := r.router.Group("books")
	// all users
	bookRoute.GET("", r.authMiddleware, r.controller.GetAll)
	bookRoute.GET("/:id", r.authMiddleware, r.controller.GetById)
	// admin only
	bookRoute.POST("", r.authAdminMiddleware, r.controller.Store)
	bookRoute.PUT("/:id", r.authAdminMiddleware, r.controller.Update)
	bookRoute.DELETE("/:id", r.authAdminMiddleware, r.controller.Delete)
}
