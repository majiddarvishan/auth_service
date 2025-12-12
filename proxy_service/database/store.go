package database

import "fmt"

// Store defines CRUD operations for all models.
type Store interface {
	// Init prepares the store (e.g. connect to DB).
	Init() error

	// User
	CreateUser(u *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(u *User) error
	DeleteUser(id uint) error
	DeleteUserByUsername(username string) error

	GetUserAndRoleByUsername(username string) (*User, error)
	UpdateUserRoleByUsername(username, roleName string) error
	
	// Multi-role support
	AddRolesToUser(userID uint, roleNames []string) error
	RemoveRolesFromUser(userID uint, roleNames []string) error
	GetUserRoles(userID uint) ([]Role, error)
	SetUserRoles(userID uint, roleNames []string) error // Replace all roles

    GetUserPhones(userName string) ([]string, error)
    AddPhoneForUser(username string, phones []string) error

	// Role
	CreateRole(r *Role) error
	GetRoleByID(id uint) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	GetAllRoles() ([]Role, error)
	GetRoles(limit, offset int) ([]Role, int64, error) // Paginated roles
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
	GetAllCustomEndpoints() ([]CustomEndpoint, error)
	UpdateCustomEndpoint(c *CustomEndpoint) error
	DeleteCustomEndpoint(id uint) error
    DeleteCustomEndpointByPath(path string) error

	// RefreshToken
	CreateRefreshToken(rt *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)
	ValidateRefreshToken(token string) (*RefreshToken, error) // Check if valid and not revoked
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID uint) error // Logout from all devices
	DeleteExpiredTokens() error // Cleanup expired tokens
}

var DB Store

// NewStore creates and initializes a Store based on kind:
// "postgres" for PGStore, "mock" for MockStore.
// dsn is passed to Init(dsn).
func NewStore(kind string) (Store, error) {
	switch kind {
	case "postgres":
		DB = NewPGStore()
	case "mock":
		DB = NewMockStore()
	default:
		return nil, fmt.Errorf("unknown store kind: %s", kind)
	}
	if err := DB.Init(); err != nil {
		return nil, err
	}
	return DB, nil
}
