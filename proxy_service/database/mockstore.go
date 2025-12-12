package database

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type MockStore struct {
	mu              sync.Mutex
	users           map[uint]*User
	roles           map[uint]*Role
	accountingRules map[uint]*AccountingRule
	customEndpoints map[uint]*CustomEndpoint
	nextID          uint
}

func NewMockStore() *MockStore {
	return &MockStore{
		users:           make(map[uint]*User),
		roles:           make(map[uint]*Role),
		accountingRules: make(map[uint]*AccountingRule),
		customEndpoints: make(map[uint]*CustomEndpoint),
		nextID:          1,
	}
}

// Init is a no-op for MockStore.
func (m *MockStore) Init() error {

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	if err != nil {
		return errors.New("Could not hash password")
	}

	adminRole := &Role{
		Name:        "admin",
		Description: "Administrator role",
	}
	adminRole.ID = m.allocateID()
	m.roles[adminRole.ID] = adminRole

	u := &User{
		Username: "admin",
		Password: string(hashedPassword),
		Roles:    []Role{*adminRole},
		Balance:  10000,
	}

	u.ID = m.allocateID()
	m.users[u.ID] = u
	return nil
}

func (m *MockStore) allocateID() uint {
	id := m.nextID
	m.nextID++
	return id
}

// User
func (m *MockStore) CreateUser(u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u.ID = m.allocateID()
	m.users[u.ID] = u
	return nil
}

func (m *MockStore) GetUserByID(id uint) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockStore) GetUserByUsername(username string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockStore) GetUserAndRoleByUsername(username string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockStore) UpdateUserRoleByUsername(username, roleName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.users {
		if u.Username == username {
			// Find the role
			for _, r := range m.roles {
				if r.Name == roleName {
					u.Roles = []Role{*r}
					return nil
				}
			}
			return errors.New("role not found")
		}
	}
	return errors.New("user not found")
}

// AddRolesToUser adds one or more roles to a user
func (m *MockStore) AddRolesToUser(userID uint, roleNames []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	user, ok := m.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	for _, roleName := range roleNames {
		found := false
		for _, r := range m.roles {
			if r.Name == roleName {
				// Check if role already exists
				exists := false
				for _, ur := range user.Roles {
					if ur.ID == r.ID {
						exists = true
						break
					}
				}
				if !exists {
					user.Roles = append(user.Roles, *r)
				}
				found = true
				break
			}
		}
		if !found {
			return errors.New("role not found: " + roleName)
		}
	}
	return nil
}

// RemoveRolesFromUser removes one or more roles from a user
func (m *MockStore) RemoveRolesFromUser(userID uint, roleNames []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	user, ok := m.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	for _, roleName := range roleNames {
		newRoles := []Role{}
		found := false
		for _, ur := range user.Roles {
			if ur.Name == roleName {
				found = true
			} else {
				newRoles = append(newRoles, ur)
			}
		}
		if !found {
			return errors.New("role not found: " + roleName)
		}
		user.Roles = newRoles
	}
	return nil
}

// GetUserRoles returns all roles for a user
func (m *MockStore) GetUserRoles(userID uint) ([]Role, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	user, ok := m.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user.Roles, nil
}

// SetUserRoles replaces all user roles with the given ones
func (m *MockStore) SetUserRoles(userID uint, roleNames []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	user, ok := m.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	newRoles := []Role{}
	for _, roleName := range roleNames {
		found := false
		for _, r := range m.roles {
			if r.Name == roleName {
				newRoles = append(newRoles, *r)
				found = true
				break
			}
		}
		if !found {
			return errors.New("role not found: " + roleName)
		}
	}
	user.Roles = newRoles
	return nil
}

func (m *MockStore) GetAllUsers() ([]User, error) {
	return nil, nil

}

func (m *MockStore) UpdateUser(u *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[u.ID]; !ok {
		return errors.New("user not found")
	}
	m.users[u.ID] = u
	return nil
}

func (m *MockStore) DeleteUser(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockStore) DeleteUserByUsername(username string) error {
	return nil
}

func (s *MockStore) GetUserPhones(userName string) ([]string, error) {
	return nil, nil
}

func (s *MockStore) AddPhoneForUser(username string, phones []string) error {
	return nil
}

// Role
func (m *MockStore) CreateRole(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r.ID = m.allocateID()
	m.roles[r.ID] = r
	return nil
}

func (m *MockStore) GetRoleByID(id uint) (*Role, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.roles[id]; ok {
		return r, nil
	}
	return nil, errors.New("role not found")
}

func (m *MockStore) GetRoleByName(name string) (*Role, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range m.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, errors.New("role not found")
}

func (m *MockStore) GetAllRoles() ([]Role, error) {
	return nil, nil
}

func (m *MockStore) UpdateRole(r *Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.roles[r.ID]; !ok {
		return errors.New("role not found")
	}
	m.roles[r.ID] = r
	return nil
}

func (m *MockStore) DeleteRole(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.roles[id]; !ok {
		return errors.New("role not found")
	}
	delete(m.roles, id)
	return nil
}

// AccountingRule
func (m *MockStore) CreateAccountingRule(a *AccountingRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	a.ID = m.allocateID()
	m.accountingRules[a.ID] = a
	return nil
}

func (m *MockStore) GetAccountingRuleByID(id uint) (*AccountingRule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if a, ok := m.accountingRules[id]; ok {
		return a, nil
	}
	return nil, errors.New("accounting rule not found")
}

func (m *MockStore) GetAccountingRuleByEndpoint(endpoint string) (*AccountingRule, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.accountingRules {
		if a.Endpoint == endpoint {
			return a, nil
		}
	}
	return nil, errors.New("accounting rule not found")
}

func (m *MockStore) UpdateAccountingRule(a *AccountingRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountingRules[a.ID]; !ok {
		return errors.New("accounting rule not found")
	}
	m.accountingRules[a.ID] = a
	return nil
}

func (m *MockStore) DeleteAccountingRule(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountingRules[id]; !ok {
		return errors.New("accounting rule not found")
	}
	delete(m.accountingRules, id)
	return nil
}

// CustomEndpoint
func (m *MockStore) CreateCustomEndpoint(c *CustomEndpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c.ID = m.allocateID()
	m.customEndpoints[c.ID] = c
	return nil
}

func (m *MockStore) GetCustomEndpointByID(id uint) (*CustomEndpoint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.customEndpoints[id]; ok {
		return c, nil
	}
	return nil, errors.New("custom endpoint not found")
}

func (m *MockStore) GetCustomEndpointByPath(path string) (*CustomEndpoint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.customEndpoints {
		if c.Path == path {
			return c, nil
		}
	}
	return nil, errors.New("custom endpoint not found")
}

func (m *MockStore) GetAllCustomEndpoints() ([]CustomEndpoint, error) {
	var endpoints []CustomEndpoint
	return endpoints, nil
}

func (m *MockStore) UpdateCustomEndpoint(c *CustomEndpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.customEndpoints[c.ID]; !ok {
		return errors.New("custom endpoint not found")
	}
	m.customEndpoints[c.ID] = c
	return nil
}

func (m *MockStore) DeleteCustomEndpoint(id uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.customEndpoints[id]; !ok {
		return errors.New("custom endpoint not found")
	}
	delete(m.customEndpoints, id)
	return nil
}

func (m *MockStore) DeleteCustomEndpointByPath(path string) error {
    m.mu.Lock()
	defer m.mu.Unlock()

    id := -1
	for _, c := range m.customEndpoints {
		if c.Path == path {
            id = int(c.ID)
            break
		}
	}

    if id == -1 {
        return errors.New("custom endpoint not found")
    }

    delete(m.customEndpoints, uint(id))

	return nil
}

