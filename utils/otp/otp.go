package otp

import (
	"crypto/rand"

	"github.com/snykk/golib_backend/config"
	gomail "gopkg.in/mail.v2"
)

const otpPayloads = "0123456789"

func SendOTP(code string, receiver string) (err error) {
	configMessage := gomail.NewMessage()
	configMessage.SetHeader("From", config.AppConfig.OTPEmail)
	configMessage.SetHeader("To", receiver)
	configMessage.SetHeader("Subject", "Verification Email")
	configMessage.SetBody("text/plain", "OTP: "+code+"\nthis code wil be expired in 5 minutes")

	dialer := gomail.NewDialer("smtp.gmail.com", 587, config.AppConfig.OTPEmail, config.AppConfig.OTPPassword)

	err = dialer.DialAndSend(configMessage)
	return
}

func GenerateCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpPayloads)
	for i := 0; i < length; i++ {
		buffer[i] = otpPayloads[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
