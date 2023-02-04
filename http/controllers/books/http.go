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

type BookController struct {
	bookUsecase    book.Usecase
	ristrettoCache cache.RistrettoCache
}

func NewBookController(bookUsecase book.Usecase, ristrettoCache cache.RistrettoCache) BookController {
	return BookController{
		bookUsecase:    bookUsecase,
		ristrettoCache: ristrettoCache,
	}
}

func (c *BookController) Store(ctx *gin.Context) {
	var bookRequest requests.BookRequest

	if err := ctx.ShouldBindJSON(&bookRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	b, statusCode, err := c.bookUsecase.Store(ctxx, bookRequest.ToDomain())
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("books")

	controllers.NewSuccessResponse(ctx, statusCode, "book inserted successfully", map[string]interface{}{
		"book": responses.FromDomain(b),
	})
}

func (c *BookController) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("books"); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, "book data fetched successfully", map[string]interface{}{
			"books": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	listOfBooks, statusCode, err := c.bookUsecase.GetAll(ctxx)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	bookResponses := responses.ToResponseList(listOfBooks)

	if bookResponses == nil {
		controllers.NewSuccessResponse(ctx, statusCode, "book data is empty", []int{})
		return
	}

	go c.ristrettoCache.Set("books", bookResponses)

	controllers.NewSuccessResponse(ctx, statusCode, "book data fetched successfully", map[string]interface{}{
		"books": bookResponses,
	})
}

func (c *BookController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("book/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
			"book": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	bookDomain, statusCode, err := c.bookUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	bookResponse := responses.FromDomain(bookDomain)

	go c.ristrettoCache.Set(fmt.Sprintf("book/%d", id), bookResponse)

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("book data with id %d fetched successfully", id), map[string]interface{}{
		"book": bookResponse,
	})
}

func (c *BookController) Update(ctx *gin.Context) {
	var bookUpdateRequest requests.BookRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&bookUpdateRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	bookDomain := bookUpdateRequest.ToDomain()
	newBook, statusCode, err := c.bookUsecase.Update(ctxx, bookDomain, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("book data with id %d updated successfully", id), map[string]interface{}{
		"book": responses.FromDomain(newBook),
	})
}

func (c *BookController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	statusCode, err := c.bookUsecase.Delete(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("books", fmt.Sprintf("book/%d", id))

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("book data with id %d deleted successfully", id), nil)
}
