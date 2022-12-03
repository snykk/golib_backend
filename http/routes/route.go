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
	Auth  map[string]string `json:"auth"`
	Users map[string]string `json:"users"`
	Books map[string]string `json:"books"`
}

func RootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Base{
		Routes: Routes{
			Auth: map[string]string{
				"Login [POST]":     "/auth/login",
				"Regis [POST]":     "/auth/regis",
				"Send OTP [POST]":  "/auth/send-otp",
				"Verif OTP [POST]": "/auth/verif-otp",
			},
			Users: map[string]string{
				"Get Users [GET] <AuthorizeJWT>":             "/users",
				"Get User [GET] <AuthorizeJWT>":              "/users/:id",
				"Get User Data It Self [GET] <AuthorizeJWT>": "/users/:id",
				"Update User [PUT] <AuthorizeJWT>":           "/users",
				"Delete User [PUT] <AuthorizeJWT>":           "/users",
			},
			Books: map[string]string{
				"Get Books [GET] <AuthorizeJWT>":              "/books",
				"Get Book [GET] <AuthorizeJWT>":               "/books/:id",
				"Create Book [POST] <AuthorizeJWT> <IsAdmin>": "/books",
				"Update Book [PUT] <AuthorizeJWT> <IsAdmin>":  "/books/:id",
				"Delete Book [PUT] <AuthorizeJWT> <IsAdmin>":  "/books/:id",
			},
		},
		Middleware: map[string]string{
			"<AuthorizeJWT>": "only user with valid token can access endpoint",
			"<IsAdmin>":      "only admin can access endpoint",
		},
		Maintainer: "Moh. Najib Fikri aka snykk github.com/snykk najibfikri13@gmail.com",
		Repository: "https://github.com/snykk/golib-backend",
	})
}
