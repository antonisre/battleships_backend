package models

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
)

// User Struct
type User struct {
	gorm.Model
	Name  string `gorm:"size:100;not null;"`
	Email string `gorm:"size:100;not null;unique;"`
}

// UserJSON struct
type UserJSON struct {
	Name  string `gorm:"size:100;not null" json:"name,omitempty"`
	Email string `gorm:"size:100;not null;unique" json:"email,omitempty"`
}

// TableName Set User's table name to be `profiles`
func (UserJSON) TableName() string {
	return "users"
}

// ValidateRegister when user registering
func (user User) ValidateRegister(db *gorm.DB) error {
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("Email is required and must be valid")
	}

	return nil
}

// Validate user
func (user User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email must be valid")
		}
		return nil
	case "forgot-password":
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Email is required and must be valid")
		}
		return nil
	default:
		return nil
	}
}

// GetUserByEmail for checking the existeence user
func (user User) GetUserByEmail(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("email = ?", user.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// Register a new user
func (user *User) Register(db *gorm.DB) (*User, error) {
	var err error
	if err := db.Debug().Create(&user).Error; err != nil {
		return nil, err
	}
	return user, err
}

// GetUsers Get all users
func (userJSON UserJSON) GetUsers(begin, limit int, db *gorm.DB) (*[]UserJSON, error) {
	users := []UserJSON{}
	if err := db.Table("users").Offset(begin).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

// GetUser Get one user
func (userJSON UserJSON) GetUser(id string, db *gorm.DB) (*UserJSON, error) {
	if err := db.Debug().Table("users").Where("id = ?", id).First(&userJSON).Error; err != nil {
		return nil, err
	}
	return &userJSON, nil
}

// CountUsers from database
func (user UserJSON) CountUsers(db *gorm.DB) (int, error) {
	var count int
	if err := db.Debug().Table("users").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
