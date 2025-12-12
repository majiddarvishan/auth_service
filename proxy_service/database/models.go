package database

import (
    "github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string  `gorm:"uniqueIndex"`
	Password  string
	Roles     []Role  `gorm:"many2many:user_roles;"` // Many-to-many relationship with roles
	Balance   float64 `gorm:"default:0"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Role struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex;not null"`
	Description string
	Users     []User `gorm:"many2many:user_roles;"` // Many-to-many relationship with users
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type AccountingRule struct {
	gorm.Model
	Endpoint string  `gorm:"uniqueIndex;not null"` // Example: "/sms", "/premium_data"
	Charge   float64 `gorm:"not null"`             // Amount to charge when accessing this endpoint
}

type CustomEndpoint struct {
	gorm.Model
	Path           string         `json:"path" gorm:"uniqueIndex;not null"`                   // e.g., "/sms/*path"
	Method         string         `json:"method" gorm:"default:'ANY'"`                        // HTTP Method ("GET", "POST", etc. or ANY)
	Endpoints      pq.StringArray `json:"endpoints" gorm:"type:text[];not null;default:'{}'"` // Target endpoints
	NeedAccounting bool           `json:"needAccounting" gorm:"default:false"`                // Flag: true if route requires accounting check
	Enabled        bool           `gorm:"default:true"`
}

type Phone struct {
  gorm.Model
  Number string `gorm:"not null"`
  UserID uint   `gorm:"index;not null"`
}

// RefreshToken stores refresh tokens for users
type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"index;not null"` // Foreign key to User
	Token     string    `gorm:"uniqueIndex;not null"` // The refresh token itself
	ExpiresAt int64     `gorm:"index;not null"` // Unix timestamp when token expires
	Revoked   bool      `gorm:"default:false"`  // For token revocation
	User      User      `gorm:"foreignKey:UserID"`
}