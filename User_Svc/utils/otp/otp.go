package otp

import (
	"fmt"
	"net/smtp"
	"time"

	"math/rand"

	"github.com/spf13/viper"
)

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendOtp(UserEmail, UserName string) (string, error) {
	from := viper.GetString("EMAIL")
	password := viper.GetString("APP_PASSWORD")
	fmt.Println("from : ", from, "pass:", password)
	otp := generateOtp()
	subject := "Just One More Step! Verify Your Email to Secure Your Account"
	link := "http://localhost:8081/user/verify/?otp=" + otp + "&email=" + UserEmail
	body := fmt.Sprintf(
		"Hi %v,\n\n"+
			"We're excited to have you onboard at Gig Space! 🚀\n\n"+
			"To complete your account setup and unlock everything we have to offer, "+
			"please take a moment to verify your email address by clicking the link below:\n\n"+
			"%v\n\n"+
			"🔒 Why verify?\n"+
			"- Protect your account from unauthorized access.\n"+
			"- Ensure you receive important updates and notifications.\n\n"+
			"If you didn’t sign up for an account, no worries—you can simply ignore this email.\n\n"+
			"Looking forward to seeing you around!\n\n"+
			"Cheers,\n"+
			"The Gig Space Team\n\n"+
			"P.S. This link will expire soon, so don’t wait too long!", UserName, link)

	msg := "From: " + from + "\n" +
		"To: " + UserEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return otp, smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{UserEmail}, []byte(msg))
}

func ForgotOtp(UserEmail, UserName string) (string, error) {
	from := viper.GetString("Email")
	password := viper.GetString("AppPassword")

	otp := generateOtp()

	subject := "Password Reset OTP"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"We received a request to reset your password. Use the OTP below to proceed:\n\n"+
			"OTP: %s\n\n"+
			"If you did not request a password reset, please ignore this email.\n\n"+
			"Best regards,\n\nGig Space Team",
		UserName, otp)

	msg := "From: " + from + "\n" +
		"To: " + UserEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return otp, smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{UserEmail}, []byte(msg))
}
