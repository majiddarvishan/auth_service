package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"auth_service/config" // Replace "auth_service" with your module name.
)

var DB *gorm.DB

// User represents a user in the database.
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Role     string // E.g., "admin", "user", etc.
}

// InitDB connects to the database and runs any necessary migrations.
func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the User model.
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatal("Failed to auto migrate database:", err)
	}
}
