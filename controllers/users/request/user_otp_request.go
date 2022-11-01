package request

type UserSendOTP struct {
	Email string `json:"email" validate:"required"`
}

type UserVerifOTP struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}
