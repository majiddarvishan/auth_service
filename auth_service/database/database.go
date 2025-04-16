package database

import (
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "auth_service/config"
)

var DB *gorm.DB

// User represents a user in the database.
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex"`
    Password string
    Role     string // E.g., "admin", "user", etc.
}

// Role represents a role definition in the system.
type Role struct {
    gorm.Model
    Name        string `gorm:"uniqueIndex;not null"`
    Description string
}

// InitDB connects to the database and migrates the schema.
func InitDB() {
    var err error
    DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate the User and Role models.
    if err := DB.AutoMigrate(&User{}, &Role{}); err != nil {
        log.Fatal("Failed to auto migrate database:", err)
    }
}
