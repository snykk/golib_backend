package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Base struct {
	Routes     Routes            `json:"routes"`
	Middleware map[string]string `json:"middleware"`
	Maintainer string            `json:"maintainer"`
	Repository string            `json:"repository"`
}

type Routes struct {
	Auth    map[string]string `json:"auth"`
	Users   map[string]string `json:"users"`
	Books   map[string]string `json:"books"`
	Reviews map[string]string `json:"reviews"`
}

func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Base{
		Routes: Routes{
			Auth: map[string]string{
				"login [POST]":     "/auth/login",
				"regis [POST]":     "/auth/regis",
				"send OTP [POST]":  "/auth/send-otp",
				"verif OTP [POST]": "/auth/verif-otp",
			},
			Users: map[string]string{
				"get all users [GET] <CommonTokenJWT>":    "/users",
				"get user by id [GET] <CommonTokenJWT>":   "/users/:id",
				"get user data [GET] <CommonTokenJWT>":    "/users/me",
				"update user data [PUT] <CommonTokenJWT>": "/users",
				"delete user [DELETE] <CommonTokenJWT>":   "/users",
				"change email [POST] <CommonTokenJWT>":    "/users/change-email",
				"change password [POST] <CommonTokenJWT>": "/users/change-password",
			},
			Books: map[string]string{
				"get all books [GET] <CommonTokenJWT>":  "/books",
				"get book by id [GET] <CommonTokenJWT>": "/books/:id",
				"create book [POST] <AdminTokenJWT>":    "/books",
				"update book [PUT] <AdminTokenJWT>":     "/books/:id",
				"delete book [DELETE] <AdminTokenJWT>":  "/books/:id",
			},
			Reviews: map[string]string{
				"get all reviews [GET] <CommonTokenJWT>":       "/reviews",
				"get review by id [GET] <CommonTokenJWT>":      "/reviews/:id",
				"get review by book id [GET] <CommonTokenJWT>": "/reviews/book/:id",
				"get review by user id [GET] <CommonTokenJWT>": "/reviews/user/:id",
				"create review [POST] <CommonTokenJWT>":        "/reviews",
				"update review [PUT] <CommonTokenJWT>":         "/reviews/:id",
				"delete review [DELETE] <CommonTokenJWT>":      "/reviews/:id",
			},
		},
		Middleware: map[string]string{
			"<CommonTokenJWT>": "user with valid basic token can access endpoint",
			"<AdminTokenJWT>":  "only user with valid admin token can access endpoint",
		},
		Maintainer: "Moh. Najib Fikri aka snykk | github.com/snykk | najibfikri13@gmail.com",
		Repository: "https://github.com/snykk/golib-backend",
	})
}
