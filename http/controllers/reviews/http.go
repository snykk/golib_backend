package reviews

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/datasources/cache"
	"github.com/snykk/golib_backend/domains/reviews"
	"github.com/snykk/golib_backend/helpers"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/reviews/requests"
	"github.com/snykk/golib_backend/http/controllers/reviews/responses"
	"github.com/snykk/golib_backend/http/token"
)

type ReviewController struct {
	reviewUsecase  reviews.Usecase
	ristrettoCache cache.RistrettoCache
}

func NewReviewController(reviewUsecase reviews.Usecase, ristrettoCache cache.RistrettoCache) ReviewController {
	return ReviewController{
		reviewUsecase:  reviewUsecase,
		ristrettoCache: ristrettoCache,
	}
}

func (c *ReviewController) Store(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	var reviewRequest requests.ReviewRequest
	if err := ctx.ShouldBindJSON(&reviewRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := helpers.IsRatingValid(reviewRequest.Rating); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	reviewDom := reviewRequest.ToDomain()

	userReview, _, _ := c.reviewUsecase.GetUserReview(ctxx, reviewDom.BookId, userClaims.UserID)
	if userReview.ID != 0 {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, "user already make a review")
		return
	}
	review, statusCode, err := c.reviewUsecase.Store(ctxx, reviewDom, userClaims.UserID)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("reviews", "users", fmt.Sprintf("user/%d", userClaims.UserID), "books", fmt.Sprintf("book/%d", reviewRequest.BookId))

	controllers.NewSuccessResponse(ctx, statusCode, "review created successfully", gin.H{
		"reviews": responses.FromDomain(review),
	})
}

func (c *ReviewController) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("reviews"); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, "review data fetched successfully", map[string]interface{}{
			"reviews": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	listOfReviews, statusCode, err := c.reviewUsecase.GetAll(ctxx)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	reviews := responses.ToResponseList(listOfReviews)

	if reviews == nil {
		controllers.NewSuccessResponse(ctx, statusCode, "review data is empty", []int{})
		return
	}

	go c.ristrettoCache.Set("reviews", reviews)

	controllers.NewSuccessResponse(ctx, statusCode, "review data fetched successfully", map[string]interface{}{
		"reviews": reviews,
	})
}

func (c *ReviewController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("review/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, http.StatusOK, fmt.Sprintf("review data with id %d fetched successfully", id), map[string]interface{}{
			"review": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	bookDomain, statusCode, err := c.reviewUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	bookResponse := responses.FromDomain(bookDomain)

	go c.ristrettoCache.Set(fmt.Sprintf("review/%d", id), bookResponse)

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("review data with id %d fetched successfully", id), map[string]interface{}{
		"review": bookResponse,
	})
}

func (c *ReviewController) GetByBookId(ctx *gin.Context) {
	bookId, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	reviewsDomain, statusCode, err := c.reviewUsecase.GetByBookId(ctxx, bookId)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	reviews := responses.ToResponseList(reviewsDomain)

	if reviews == nil {
		controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("review data with book id %d is empty", bookId), []int{})
		return
	}

	controllers.NewSuccessResponse(ctx, http.StatusOK, fmt.Sprintf("review data with book id %d fetched successfully", bookId), map[string]interface{}{
		"review": reviews,
	})
}

func (c *ReviewController) GetByUserid(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("id"))
	ctxx := ctx.Request.Context()

	reviewsDomain, statusCode, err := c.reviewUsecase.GetByUserId(ctxx, userId)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	reviews := responses.ToResponseList(reviewsDomain)

	if reviews == nil {
		controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("review data with user id %d is empty", userId), []int{})
		return
	}

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("review data with user id %d fetched successfully", userId), map[string]interface{}{
		"review": reviews,
	})
}

func (c *ReviewController) Update(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	reviewId, _ := strconv.Atoi(ctx.Param("id"))
	var reviewRequest requests.ReviewRequest
	if err := ctx.ShouldBindJSON(&reviewRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := helpers.IsRatingValid(reviewRequest.Rating); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	reviewDom := reviewRequest.ToDomain()
	review, statusCode, err := c.reviewUsecase.Update(ctxx, reviewDom, userClaims.UserID, reviewId)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("reviews", fmt.Sprintf("review/%d", review.ID), "users", fmt.Sprintf("user/%d", userClaims.UserID), "books", fmt.Sprintf("book/%d", reviewRequest.BookId))

	controllers.NewSuccessResponse(ctx, statusCode, "review updated successfully", gin.H{
		"reviews": responses.FromDomain(review),
	})
}

func (c *ReviewController) Delete(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	reviewid, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	bookId, statusCode, err := c.reviewUsecase.Delete(ctxx, userClaims.UserID, reviewid)
	if err != nil {
		controllers.NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("reviews", fmt.Sprintf("review/%d", reviewid), "books", fmt.Sprintf("book/%d", userClaims.UserID), "books", fmt.Sprintf("book/%d", bookId))

	controllers.NewSuccessResponse(ctx, statusCode, fmt.Sprintf("review data with id %d deleted successfully", reviewid), nil)
}
