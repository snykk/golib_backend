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
	book "github.com/snykk/golib_backend/usecase/books"
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

func (controller BookController) Store(c *gin.Context) {
	var book requests.BookRequest
	var err error
	err = c.Bind(&book)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	isValid, err := isInsertedBookValid(&book)
	if !isValid {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()

	b, err := controller.BookUsecase.Store(ctx, book.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, "book inserted successfully", map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (controller BookController) GetAll(c *gin.Context) {
	b, err := controller.BookUsecase.GetAll()

	var books []responses.BookResponse

	for _, val := range b {
		books = append(books, responses.FromDomain(val))
	}

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if books == nil {
		controllers.NewSuccessResponse(c, "book data is empty", []int{})
		return
	}

	controllers.NewSuccessResponse(c, "book data fetched successfully", books)
	return
}

func (controller BookController) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	ctx := c.Request.Context()
	result, err := controller.BookUsecase.GetById(ctx, id)
	if err != nil {
		controllers.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
		"book": responses.FromDomain(result),
	})
}

func (controller BookController) Update(c *gin.Context) {
	var book requests.BookRequest
	id, _ := strconv.Atoi(c.Param("id"))
	ctx := c.Request.Context()

	c.Bind(&book)

	isValid, err := isInsertedBookValid(&book)

	if !isValid {
		controllers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bookDomain := book.ToDomain()
	bookDomain.ID = id
	b, err := controller.BookUsecase.Update(ctx, bookDomain)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("book data with id %d updated successfully", id), map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (controller BookController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	ctx := c.Request.Context()
	err := controller.BookUsecase.Delete(ctx, id)

	if err != nil {
		controllers.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	controllers.NewSuccessResponse(c, fmt.Sprintf("book data with id %d deleted successfully", id), nil)
}
