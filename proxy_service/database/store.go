package database

import "gorm.io/gorm"

// Store defines CRUD operations for all models.
type Store interface {
    // User
    CreateUser(u *User) error
    GetUserByID(id uint) (*User, error)
    GetUserByUsername(username string) (*User, error)
    UpdateUser(u *User) error
    DeleteUser(id uint) error

    // Role
    CreateRole(r *Role) error
    GetRoleByID(id uint) (*Role, error)
    GetRoleByName(name string) (*Role, error)
    UpdateRole(r *Role) error
    DeleteRole(id uint) error

    // AccountingRule
    CreateAccountingRule(a *AccountingRule) error
    GetAccountingRuleByID(id uint) (*AccountingRule, error)
    GetAccountingRuleByEndpoint(endpoint string) (*AccountingRule, error)
    UpdateAccountingRule(a *AccountingRule) error
    DeleteAccountingRule(id uint) error

    // CustomEndpoint
    CreateCustomEndpoint(c *CustomEndpoint) error
    GetCustomEndpointByID(id uint) (*CustomEndpoint, error)
    GetCustomEndpointByPath(path string) (*CustomEndpoint, error)
    UpdateCustomEndpoint(c *CustomEndpoint) error
    DeleteCustomEndpoint(id uint) error
}

// NewPGStore creates a GORM-based Store.
func NewPGStore(db *gorm.DB) Store {
    return &PGStore{db: db}
}