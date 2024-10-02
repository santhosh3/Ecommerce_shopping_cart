package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)


func SendOTP(smtpPORT, otp int, smtpHost, emailHOST, receiverEMAIL, passwordHOST string) error {

	// Create an HTML template for the email body
	htmlBody := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<body>
				<h2>Password Reset Request</h2>
				<p>Use the following OTP to change your password:</p>
				<h2>%06d</h2>
				<p>If you did not request a password reset, please ignore this email.</p>
			</body>
			</html>
		`, otp)

	// Create a new dialer
	dialer := gomail.NewDialer(smtpHost, smtpPORT, emailHOST, passwordHOST)

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", emailHOST)
	msg.SetHeader("To", receiverEMAIL) // Replace with recipient's email
	msg.SetHeader("Subject", "OTP for reset password")
	msg.SetBody("text/html", htmlBody)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}