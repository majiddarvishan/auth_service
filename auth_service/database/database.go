package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth_service/config" // update "auth_service" to your module name
)

var DB *gorm.DB

// User model defines a user in the database.
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     string
}

// InitDB opens the database connection and migrates the schema.
func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatal("Failed to migrate User schema:", err)
	}
}
