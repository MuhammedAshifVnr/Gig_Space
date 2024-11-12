package helper

import (
	"fmt"
	"net/smtp"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
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

	var subject, message string

	switch topic {
	case "Fail":
		subject = "Refund Notification - Wallet Required for Refund"
		message = fmt.Sprintf(
			"Dear User,\n\n"+
				"We wanted to inform you that a recent order was canceled, and a refund of $%d is due to you. "+
				"Our system issues refunds to your account wallet. However, it appears you currently do not have an active wallet associated with your account.\n\n"+
				"To proceed with the refund, please set up a wallet in your account settings. Once your wallet is active, we will process the refund at the earliest.\n\n"+
				"If you need assistance, feel free to reach out to our support team.\n\n"+
				"Thank you for your attention.\n\n"+
				"Best regards,\n"+
				"Gig Space Team", Amount)
	case "Done":
		subject = "Payment Credited to Your Wallet"
		message = fmt.Sprintf(
			"Dear Freelancer,\n\n"+
				"We are pleased to inform you that a payment of $%d has been credited to your wallet for successfully completing the order with ID: %s.\n\n"+
				"You can view your updated wallet balance in your account. If you have any questions or need further assistance, please reach out to our support team.\n\n"+
				"Thank you for your hard work and dedication.\n\n"+
				"Best regards,\n"+
				"Gig Space Team", Amount, orderID)
	default:
		subject = "Notification"
		message = "No specific details available for this notification."
	}

	return subject, message
}

func ForgotMsgCreater(topic, otp string) (string, string) {
	var subject, message string
	switch topic {
	case "Wallet":
		subject = "Password Reset OTP"
		message = fmt.Sprintf(
			"Hello User,\n\n"+
				"We received a request to reset your password. Use the OTP below to proceed:\n\n"+
				"OTP: %s\n\n"+
				"If you did not request a password reset, please ignore this email.\n\n"+
				"Best regards,\n\nGig Space Team", otp)
	}
	return subject, message
}

func OfflineMessage(senderName string) (string, string) {
	subject := "New Message Notification"

	message := fmt.Sprintf(
		"Hello User,\n\n"+
			"You have received a new message from %s on Gig Space.\n\n"+
			"To view and reply to the message, please log in to your account.\n\n"+
			"Best regards,\n\nGig Space Team", senderName)

	return subject, message
}

func OrderMessages(Order *proto.OrderDetail, topic string) (string, string, string) {
	var subject, message, email string
	switch topic {
	case "OrderReceived":
		email, _ = GetUserEmail(uint(Order.FrelancerId))
		subject = "You Have a New Order - Action Required"
		message = fmt.Sprintf(
			"Hello,\n\n"+
				"We’re excited to inform you that you have received a new order! Here are the order details:\n\n"+
				"Order ID: %s\n\n"+
				"Please take a moment to review the order. Once you’ve checked all the requirements, you can choose to either accept or reject the order based on your availability and preferences.\n\n"+
				"**To Accept the Order**: Log in to your account, review the order details, and click 'Accept'.\n"+
				"**To Reject the Order**: If you are unable to fulfill this order, you may click 'Reject' to notify the client.\n\n"+
				"Note: Accepting the order means you agree to meet the agreed delivery timeline and project scope.\n\n"+
				"Thank you for using Gig Space. We're here to support your success every step of the way.\n\n"+
				"Best regards,\n"+
				"The Gig Space Team", Order.OrderId)
	}
	return subject, message, email
}
