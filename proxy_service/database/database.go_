package database

import (
	"fmt"
	"log"

	"auth_service/config"

	"gorm.io/driver/postgres"
    "github.com/lib/pq"
	"gorm.io/gorm"
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
	Username string `gorm:"uniqueIndex"`
	Password string
    RoleID   uint // Foreign key field
	Role     Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Association

	// Role     string // Optional if you want a backup or quick check.
	// RoleID   uint     // Foreign key to the Role table.
	// RoleInfo Role     `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Balance float64 `gorm:"default:0"`
}

// AccountingRule defines which endpoints require a balance check and their charge.
type AccountingRule struct {
	gorm.Model
	Endpoint string  `gorm:"uniqueIndex;not null"` // Example: "/sms", "/premium_data"
	Charge   float64 `gorm:"not null"`             // Amount to charge when accessing this endpoint
}

// CustomEndpoint represents a user-defined route configuration.
type CustomEndpoint struct {
	gorm.Model
	Path   string `json:"path" gorm:"uniqueIndex;not null"` // e.g., "/sms/*path"
	Method string `json:"method" gorm:"default:'ANY'"`      // HTTP Method ("GET", "POST", etc. or ANY)
	Endpoints      pq.StringArray `json:"endpoints" gorm:"type:text[];not null;default:'{}'"` // Target endpoints
	NeedAccounting bool           `json:"needAccounting" gorm:"default:false"`                // Flag: true if route requires accounting check
	Enabled        bool           `gorm:"default:true"`
}

// InitDB initializes the database and performs migrations.
func InitDB() {
	var err error

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable  TimeZone=UTC",
		config.DatabaseHost, config.DatabasePort, config.DatabaseUserName, config.DatabasePassword, config.DatabaseName)

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models.
	if err := DB.AutoMigrate(&User{}, &Role{}, &AccountingRule{}, &CustomEndpoint{}); err != nil {
		log.Fatal("Failed to auto migrate database:", err)
	}
}
