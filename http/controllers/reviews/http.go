package reviews

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/snykk/golib_backend/constants"
	"github.com/snykk/golib_backend/domains/reviews"
	"github.com/snykk/golib_backend/http/controllers"
	"github.com/snykk/golib_backend/http/controllers/reviews/requests"
	"github.com/snykk/golib_backend/http/controllers/reviews/responses"
	"github.com/snykk/golib_backend/packages/cache"
	"github.com/snykk/golib_backend/packages/token"
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

func (c *ReviewController) AddReview(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	var reviewRequest requests.ReviewRequest
	if err := ctx.ShouldBindJSON(&reviewRequest); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := isValidRating(reviewRequest.Rating); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	reviewDom := reviewRequest.ToDomain()
	review, err := c.reviewUsecase.Store(ctxx, reviewDom, userClaims.UserID)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	controllers.NewSuccessResponse(ctx, "review created successfully", gin.H{
		"reviews": responses.FromDomain(review),
	})
}

func (c *ReviewController) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("reviews"); val != nil {
		controllers.NewSuccessResponse(ctx, "review data fetched successfully", map[string]interface{}{
			"reviews": val,
		})
		return
	}

	listOfReviews, err := c.reviewUsecase.GetAll()
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	reviews := responses.ToResponseList(listOfReviews)

	if reviews == nil {
		controllers.NewSuccessResponse(ctx, "review data is empty", []int{})
		return
	}

	go c.ristrettoCache.Set("reviews", reviews)

	controllers.NewSuccessResponse(ctx, "review data fetched successfully", map[string]interface{}{
		"books": reviews,
	})
}

func (c *ReviewController) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("review/%d", id)); val != nil {
		controllers.NewSuccessResponse(ctx, fmt.Sprintf("review data with id %d fetched successfully", id), map[string]interface{}{
			"review": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	bookDomain, err := c.reviewUsecase.GetById(ctxx, id)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	bookResponse := responses.FromDomain(bookDomain)

	go c.ristrettoCache.Set(fmt.Sprintf("review/%d", id), bookResponse)

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("review data with id %d fetched successfully", id), map[string]interface{}{
		"review": bookResponse,
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

	if err := isValidRating(reviewRequest.Rating); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	reviewDom := reviewRequest.ToDomain()
	review, err := c.reviewUsecase.Update(ctxx, reviewDom, userClaims.UserID, reviewId)
	if err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("reviews", fmt.Sprintf("review/%d", review.ID))

	controllers.NewSuccessResponse(ctx, "review created successfully", gin.H{
		"reviews": responses.FromDomain(review),
	})
}

func (c *ReviewController) Delete(ctx *gin.Context) {
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(token.JwtCustomClaim)
	reviewid, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	if err := c.reviewUsecase.Delete(ctxx, userClaims.UserID, reviewid); err != nil {
		controllers.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go c.ristrettoCache.Del("reviews", fmt.Sprintf("review/%d", reviewid))

	controllers.NewSuccessResponse(ctx, fmt.Sprintf("review data with id %d deleted successfully", reviewid), nil)
}

func isValidRating(rating int) error {
	if rating < 1 || rating > 10 {
		return errors.New("the rating must be in the range 1 - 10")
	}
	return nil
}
