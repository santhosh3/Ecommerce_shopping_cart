package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/gomail.v2"
	"gorm.io/datatypes"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	defer r.Body.Close() // Ensure the body is closed after reading

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Optional: Disallow unknown fields to avoid silent errors

	err := decoder.Decode(payload)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(req *http.Request) string {
	tokenAuth := req.Header.Get("Authorization")
	tokenQuery := req.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}
	return ""
}

func GenerateOTP() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000 // Generates a 6-digit OTP
}

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

func ConvertStringToFloat(f string) float64 {
	if s, err := strconv.ParseFloat(f, 32); err == nil {
		return s
	}
	return 0
}

func ConvertStringToBool(f string) bool {
	boolValue, err := strconv.ParseBool(f)
	if err != nil {
		return false
	}
	return boolValue
}

func ConvertStringToArray(input string) (datatypes.JSON) {
	// Split the string by commas and trim spaces
	elements := strings.Split(input, ",")
	for i := range elements {
		elements[i] = strings.TrimSpace(elements[i])
	}

	// Convert to JSON
	jsonData, err := json.Marshal(elements)
	if err != nil {
		return jsonData
	}
	return jsonData
}