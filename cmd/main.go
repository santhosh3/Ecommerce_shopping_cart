// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/santhosh3/ECOM/Config"
// 	"github.com/santhosh3/ECOM/cmd/api"
// 	"github.com/santhosh3/ECOM/database"
// 	"github.com/santhosh3/ECOM/models"
// 	"gorm.io/gorm"
// )

// func main() {
// 	//taking psqlString from ENV
// 	connectionString := config.Envs.PostgresString
// 	if len(connectionString) == 0 {
// 		log.Fatal("POSTGRES_SQL is not set in .env file")
// 	}

// 	//connecting to postgres DB
// 	db, err := database.NewPSQLStorage(connectionString)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// doing migrations
// 	models.DBMigrations(db)

// 	//checking DB connections
// 	initStorage(db)

// 	//running API server
// 	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
// 	if err := server.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func initStorage(db *gorm.DB) {
// 	sqlDB, err := db.DB() // Get the underlying sql.DB object from the GORM DB object
// 	if err != nil {
// 		log.Fatal("Failed to get database handle:", err)
// 	}

// 	// Ping the database to check if it's reachable
// 	err = sqlDB.Ping()
// 	if err != nil {
// 		log.Fatal("Failed to connect to the database:", err)
// 	}

// 	log.Println("DB: connected successfully")
// }

package main

import (
	"fmt"
	"time"
	"math/rand"
	"gopkg.in/gomail.v2"
)

func generateOTP() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000 // Generates a 6-digit OTP
}

func main() {
	// Load your credentials from environment variables for security
	smtpHost := "smtp.office365.com"
	smtpPort := 587
	smtpUser := "santhoshchinna109@outlook.com" // Replace with your Outlook email
	smtpPass := "Chinna@123"                    // Ensure this is set in your environment variables

	// Check if password environment variable is set
	if smtpPass == "" {
		fmt.Println("OUTLOOK_SMTP_PASS environment variable is not set")
		return
	}

	// Generate a random OTP
	otp := generateOTP()

	// Create an HTML template for the email body
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>Use the following OTP to change your password:</p>
			<h3>%06d</h3>
			<p>If you did not request a password reset, please ignore this email.</p>
		</body>
		</html>
	`, otp)

	// Create a new dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", smtpUser)
	msg.SetHeader("To", "santhoshchinna109@gmail.com") // Replace with recipient's email
	msg.SetHeader("Subject", "Test Email from Outlook")
	msg.SetBody("text/html", htmlBody)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		fmt.Println("Failed to send email:", err)
		// Handle error appropriately
		return
	}

	fmt.Println("Email sent successfully!")
}
