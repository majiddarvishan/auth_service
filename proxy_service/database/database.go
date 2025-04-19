package database

import (
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "auth_service/config"
)

var DB *gorm.DB

type Role struct {
    gorm.Model
    Name        string `gorm:"uniqueIndex;not null"`
    Description string
}

// User represents a user with a balance for charging purposes.
type User struct {
    gorm.Model
    Username string   `gorm:"uniqueIndex"`
    Password string
    Role     string   // Optional if you want a backup or quick check.
    // RoleID   uint     // Foreign key to the Role table.
    // RoleInfo Role     `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Balance  float64  `gorm:"default:0"`
}

// AccountingRule defines which endpoints require a balance check and their charge.
type AccountingRule struct {
    gorm.Model
    Endpoint string  `gorm:"uniqueIndex;not null"` // Example: "/sms", "/premium_data"
    Charge   float64 `gorm:"not null"`             // Amount to charge when accessing this endpoint
}

// InitDB initializes the database and performs migrations.
func InitDB() {
    var err error
    DB, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate models.
    if err := DB.AutoMigrate(&User{}, &Role{}, &AccountingRule{}); err != nil {
        log.Fatal("Failed to auto migrate database:", err)
    }
}
