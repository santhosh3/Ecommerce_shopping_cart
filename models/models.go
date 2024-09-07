package models

import (
	"gorm.io/gorm"
)

// DBMigration performs the database migrations for the defined models.
func DBMigration(db *gorm.DB) {
	// Concatenating UserModel and ProductModel
	modelsToMigrate := append(UserModel, ProductModel...)

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(&model); err != nil {
			panic("Failed to migrate model: " + err.Error())
		}
	}
}