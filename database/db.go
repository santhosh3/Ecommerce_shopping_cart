package database

import (
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPSQLStorage(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBQueryTimeoutMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start timing the request
			start := time.Now()

			// Call the next handler in the chain
			next.ServeHTTP(w, r)

			// Measure the execution time
			log.Printf("Response time: %s", time.Since(start))
		})
	}
}
