package request

type UserSendOTP struct {
	Email string `json:"email" binding:"required"`
}

type UserVerifOTP struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
