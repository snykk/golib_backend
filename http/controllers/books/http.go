package books

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/datasources/cache"
	book "github.com/snykk/golib_backend/domains/books"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/books/requests"
	"github.com/snykk/golib_backend/http/controllers/books/responses"
)

type BookController interface {
	Store(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	BookUsecase    book.Usecase
	RistrettoCache cache.RistrettoCache
}

func NewBookController(bookUsecase book.Usecase, ristrettoCache cache.RistrettoCache) BookController {
	return &bookController{
		BookUsecase:    bookUsecase,
		RistrettoCache: ristrettoCache,
	}
}

func (bookC bookController) Store(ctx *gin.Context) {
	var bookRequest requests.BookRequest

	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	b, err := bookC.BookUsecase.Store(ctxx, bookRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "book inserted successfully", map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (bookC bookController) GetAll(ctx *gin.Context) {
	if val := bookC.RistrettoCache.Get("books"); val != nil {
		controllers.NewSuccessResponse(ctx, "book data fetched successfully", map[string]interface{}{
			"books": val,
		})
		return
	}

	listOfBooks, err := bookC.BookUsecase.GetAll()
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	bookResponses := responses.ToResponseList(listOfBooks)

	if bookResponses == nil {
		controllers.NewSuccessResponse(ctx, "book data is empty", []int{})
		return
	}

	go bookC.RistrettoCache.Set("books", bookResponses)

	controllers.NewSuccessResponse(ctx, "book data fetched successfully", map[string]interface{}{
		"books": bookResponses,
	})
}

func (bookC bookController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := bookC.RistrettoCache.Get(fmt.Sprintf("book/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
			"book": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	bookDomain, err := bookC.BookUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	bookResponse := responses.FromDomain(bookDomain)

	go bookC.RistrettoCache.Set(fmt.Sprintf("book/%d", id), bookResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
		"book": bookResponse,
	})
}

func (bookC bookController) Update(ctx *gin.Context) {
	var bookUpdateRequest requests.BookUpdateRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&bookUpdateRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctxx := ctx.Request.Context()
	bookDomain := bookUpdateRequest.ToDomain()
	newBook, err := bookC.BookUsecase.Update(ctxx, bookDomain, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	go bookC.RistrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d updated successfully", id), map[string]interface{}{
		"book": responses.FromDomain(newBook),
	})
}

func (bookC bookController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	if err := bookC.BookUsecase.Delete(ctxx, id); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	go bookC.RistrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d deleted successfully", id), nil)
}
