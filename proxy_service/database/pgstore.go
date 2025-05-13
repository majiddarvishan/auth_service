package database

import (
	"auth_service/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGStore struct {
	db *gorm.DB
}

// NewPGStore creates a GORM-based Store.
func NewPGStore() *PGStore {
	return &PGStore{}
}

// InitDB initializes the database and performs migrations.
func (s *PGStore) Init() error {
	var err error

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable  TimeZone=UTC",
		config.DatabaseHost, config.DatabasePort, config.DatabaseUserName, config.DatabasePassword, config.DatabaseName)

	s.db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models.
	if err := s.db.AutoMigrate(&User{}, &Role{}, &AccountingRule{}, &CustomEndpoint{}, Phone{}); err != nil {
		log.Fatal("Failed to auto migrate database:", err)
	}

	return nil
}

func (s *PGStore) CreateUser(u *User) error {
	return s.db.Create(u).Error
}

func (s *PGStore) GetUserByID(id uint) (*User, error) {
	var u User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PGStore) GetUserByUsername(username string) (*User, error) {
	var u User
	if err := s.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PGStore) GetUserAndRoleByUsername(username string) (*User, error) {
	var u User
	if err := s.db.
		Preload("Role").
		Where("username = ?", username).
		First(&u).Error; err != nil {
		return nil, fmt.Errorf("User not found")
	}

	return &u, nil
}

func (s *PGStore) UpdateUserRoleByUsername(username, roleName string) error {
	// Find user by username.
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("User not found")
	}

	// Look up the role in the database.
	var role Role
	if err := s.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return fmt.Errorf("Role not found")
	}

	// Update the user's role.
	user.RoleID = role.ID

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *PGStore) GetAllUsers() ([]User, error) {
	//  Load users and their Roles
	var users []User
	if err := s.db.
		Preload("Role").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("Failed to retrieve users")
	}

	return users, nil
}

func (s *PGStore) UpdateUser(u *User) error {
	return s.db.Save(u).Error
}

func (s *PGStore) DeleteUser(id uint) error {
	return s.db.Delete(&User{}, id).Error
}

func (s *PGStore) DeleteUserByUsername(username string) error {
	// Find the user in the database.
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("User not found")
	}

	// Soft delete the user .
	if err := s.db.Delete(&user).Error; err != nil {
		return err
	}

	// Permanently delete the user to clear the unique constraint.
	// if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
	//     c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user", "details": err.Error()})
	//     return
	// }

	return nil
}

func (s *PGStore) GetUserPhones(userName string) ([]string, error) {
	// Find user by username.
	var user User
	if err := s.db.Where("username = ?", userName).First(&user).Error; err != nil {
		return nil, fmt.Errorf("User not found")
	}

	var phones []Phone
	if err := s.db.Where("user_id = ?", user.ID).Find(&phones).Error; err != nil {
		return nil, fmt.Errorf("could not fetch phones")
	}

	// Extract just the numbers
	nums := make([]string, len(phones))
	for i, p := range phones {
		nums[i] = p.Number
	}

	return nums, nil
}

func (s *PGStore) AddPhoneForUser(userName string, phones []string) error {
	// Find user by username.
	var user User
	if err := s.db.Where("username = ?", userName).First(&user).Error; err != nil {
		return fmt.Errorf("User not found")
	}

	// batch insert each phone
	for _, num := range phones {
		phone := Phone{
			Number: num,
			UserID: user.ID,
		}
		if err := s.db.Create(&phone).Error; err != nil {
			return fmt.Errorf("could not save phones")

		}
	}

	return nil
}

// Role
func (s *PGStore) CreateRole(r *Role) error {
	return s.db.Create(r).Error
}

func (s *PGStore) GetRoleByID(id uint) (*Role, error) {
	var r Role
	if err := s.db.First(&r, id).Error; err != nil {
		return nil, err
	}
	return &r, nil
}
func (s *PGStore) GetRoleByName(name string) (*Role, error) {
	var r Role
	if err := s.db.Where("name = ?", name).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *PGStore) GetAllRoles() ([]Role, error) {
	var roles []Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *PGStore) UpdateRole(r *Role) error {
	return s.db.Save(r).Error
}

func (s *PGStore) DeleteRole(id uint) error {
	return s.db.Delete(&Role{}, id).Error
}

// AccountingRule
func (s *PGStore) CreateAccountingRule(a *AccountingRule) error { return s.db.Create(a).Error }
func (s *PGStore) GetAccountingRuleByID(id uint) (*AccountingRule, error) {
	var a AccountingRule
	if err := s.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
func (s *PGStore) GetAccountingRuleByEndpoint(endpoint string) (*AccountingRule, error) {
	var a AccountingRule
	if err := s.db.Where("endpoint = ?", endpoint).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *PGStore) UpdateAccountingRule(a *AccountingRule) error {
	return s.db.Save(a).Error
}

func (s *PGStore) DeleteAccountingRule(id uint) error {
	return s.db.Delete(&AccountingRule{}, id).Error
}

// CustomEndpoint
func (s *PGStore) CreateCustomEndpoint(c *CustomEndpoint) error {
	return s.db.Create(c).Error
}

func (s *PGStore) GetCustomEndpointByID(id uint) (*CustomEndpoint, error) {
	var c CustomEndpoint
	if err := s.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *PGStore) GetCustomEndpointByPath(path string) (*CustomEndpoint, error) {
	var c CustomEndpoint
	if err := s.db.Where("path = ?", path).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *PGStore) GetAllCustomEndpoints() ([]CustomEndpoint, error) {
	var endpoints []CustomEndpoint
	if err := s.db.Where("enabled = ?", true).Find(&endpoints).Error; err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (s *PGStore) UpdateCustomEndpoint(c *CustomEndpoint) error { return s.db.Save(c).Error }

func (s *PGStore) DeleteCustomEndpoint(id uint) error {
	return s.db.Delete(&CustomEndpoint{}, id).Error
}
