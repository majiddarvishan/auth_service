package database

import (
    "github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	RoleID   uint    // Foreign key field
	Role     Role    `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Balance  float64 `gorm:"default:0"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
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
