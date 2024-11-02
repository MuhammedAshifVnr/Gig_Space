package helper

import (
	"fmt"
	"net/smtp"

	"github.com/spf13/viper"
)

func SendEmailNotification(UserEmail, Subject, Message string) error {
	from := viper.GetString("Email")
	password := viper.GetString("AppPassword")

	msg := "From: " + from + "\n" +
		"To: " + UserEmail + "\n" +
		"Subject: " + Subject + "\n\n" +
		Message

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{UserEmail}, []byte(msg))
}

func MessageCreater(topic, orderID string, Amount int) (string, string) {
	if topic == "Fail" {
		subject := "Refund Notification - Wallet Required for Refund"
		message := fmt.Sprintf(
			"Dear User,\n\n"+
				"We wanted to inform you that a recent order was canceled, and a refund of $%d is due to you. Our system issues refunds to your account wallet. However, it appears you currently do not have an active wallet associated with your account.\n\n"+
				"To proceed with the refund, please set up a wallet in your account settings. Once your wallet is active, we will process the refund at the earliest.\n\n"+
				"If you need assistance, feel free to reach out to our support team.\n\n"+
				"Thank you for your attention.\n\n"+
				"Best regards,\n"+
				"Gig Space Team",
			Amount)
		return subject, message
	} else if topic == "Done" {
		subject := "Payment Credited to Your Wallet"
message := fmt.Sprintf(
	"Dear Freelancer,\n\n"+
		"We are pleased to inform you that a payment of $%d has been credited to your wallet for successfully completing the order with ID: %s.\n\n"+
		"You can view your updated wallet balance in your account. If you have any questions or need further assistance, please reach out to our support team.\n\n"+
		"Thank you for your hard work and dedication.\n\n"+
		"Best regards,\n"+
		"Gig Space Team",
	Amount, orderID)
return subject, message
	}
	return "Test Subject", "Test Message"
}
