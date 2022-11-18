package books

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/snykk/golib_backend/controllers"
	"github.com/snykk/golib_backend/controllers/books/requests"
	"github.com/snykk/golib_backend/controllers/books/responses"
	"github.com/snykk/golib_backend/datasources/cache"
	book "github.com/snykk/golib_backend/usecases/books"
)

type BookController struct {
	BookUsecase    book.Usecase
	RistrettoCache cache.RistrettoCache
}

func NewBookController(bookUsecase book.Usecase, ristrettoCache cache.RistrettoCache) *BookController {
	return &BookController{
		BookUsecase:    bookUsecase,
		RistrettoCache: ristrettoCache,
	}
}

func isInsertedBookValid(request *requests.BookRequest) (bool, error) {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (bookController BookController) Store(ctx *gin.Context) {
	var bookRequest requests.BookRequest

	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if isValid, err := isInsertedBookValid(&bookRequest); !isValid {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	b, err := bookController.BookUsecase.Store(ctxx, bookRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "book inserted successfully", map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (bookController BookController) GetAll(ctx *gin.Context) {
	if val := bookController.RistrettoCache.Get("books"); val != nil {
		controllers.NewSuccessResponse(ctx, "book data fetched successfully", map[string]interface{}{
			"books": val,
		})
		return
	}

	listOfBooks, err := bookController.BookUsecase.GetAll()
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	bookResponses := responses.ToResponseList(listOfBooks)

	if bookResponses == nil {
		controllers.NewSuccessResponse(ctx, "book data is empty", []int{})
		return
	}

	bookController.RistrettoCache.Set("books", bookResponses)

	controllers.NewSuccessResponse(ctx, "book data fetched successfully", map[string]interface{}{
		"books": bookResponses,
	})
}

func (bookController BookController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := bookController.RistrettoCache.Get(fmt.Sprintf("book/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
			"book": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	bookDomain, err := bookController.BookUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	bookResponse := responses.FromDomain(bookDomain)

	bookController.RistrettoCache.Set(fmt.Sprintf("book/%d", id), bookResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
		"book": bookResponse,
	})
}

func (bookController BookController) Update(ctx *gin.Context) {
	var bookUpdateRequest requests.BookUpdateRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&bookUpdateRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctxx := ctx.Request.Context()
	bookDomain := bookUpdateRequest.ToDomain()
	newBook, err := bookController.BookUsecase.Update(ctxx, bookDomain, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	bookController.RistrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d updated successfully", id), map[string]interface{}{
		"book": responses.FromDomain(newBook),
	})
}

func (bookController BookController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	if err := bookController.BookUsecase.Delete(ctxx, id); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	bookController.RistrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d deleted successfully", id), nil)
}
