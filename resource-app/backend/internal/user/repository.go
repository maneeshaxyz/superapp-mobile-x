package user

import (
	"gorm.io/gorm"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	GetUsers() ([]User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUserRole(userID string, role Role) error
}

// GormUserRepository is an implementation of UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetUsers() ([]User, error) {
	var users []User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *GormUserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserByID(id string) (*User, error) {
	var user User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) UpdateUserRole(userID string, role Role) error {
	return r.db.Model(&User{}).Where("id = ?", userID).Update("role", role).Error
}
