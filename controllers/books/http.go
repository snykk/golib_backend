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
	book "github.com/snykk/golib_backend/usecases/books"
)

type BookController struct {
	BookUsecase book.Usecase
}

func NewBookController(bookUsecase book.Usecase) *BookController {
	return &BookController{
		bookUsecase,
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

func (controller BookController) Store(ctx *gin.Context) {
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
	b, err := controller.BookUsecase.Store(ctxx, bookRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "book inserted successfully", map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (controller BookController) GetAll(ctx *gin.Context) {
	listOfBooks, err := controller.BookUsecase.GetAll()
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var books []responses.BookResponse
	for _, val := range listOfBooks {
		books = append(books, responses.FromDomain(val))
	}

	if books == nil {
		controllers.NewSuccessResponse(ctx, "book data is empty", []int{})
		return
	}

	controllers.NewSuccessResponse(ctx, "book data fetched successfully", books)
}

func (controller BookController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	bookDomain, err := controller.BookUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
		"book": responses.FromDomain(bookDomain),
	})
}

func (controller BookController) Update(ctx *gin.Context) {
	var bookUpdateRequest requests.BookUpdateRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&bookUpdateRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	ctxx := ctx.Request.Context()
	bookDomain := bookUpdateRequest.ToDomain()
	newBook, err := controller.BookUsecase.Update(ctxx, bookDomain, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d updated successfully", id), map[string]interface{}{
		"book": responses.FromDomain(newBook),
	})
}

func (controller BookController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	if err := controller.BookUsecase.Delete(ctxx, id); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("book data with id %d deleted successfully", id), nil)
}
