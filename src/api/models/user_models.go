package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/quattad/sudokubuddy-backend/src/api/security"
)

// User is a struct that defines the fields in the DB
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:20;not null;unique" json:"username"`
	Email     string    `gorm:"size:50;not null;unique" json:"email"`
	FirstName string    `gorm:"size:20;not null;" json:"first_name"`
	LastName  string    `gorm:"size:20;not null;" json:"last_name"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:current_timestamp()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp()" json:"updated_at"`
	Puzzles   []Puzzle  `gorm:"foreignkey:UserID" json:"puzzles"`
}

// BeforeSave generates hashed password using bcrypt from plaintext password provided in field User.Password
func (user *User) BeforeSave() error {
	hashedPassword, err := security.SecurityService.Hash(user.Password)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

// PrepareUser removes whitespaces from user fields
func (user *User) PrepareUser() *User {
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}

// ValidateUser sanitizes the user fields depending on the request type
func (user *User) ValidateUser(method string) error {
	var err error

	switch method {
	case "update":
		if user.Username == "" {
			return errors.New("User must have defined property 'username'")
		}

		if user.Email == "" {
			return errors.New("User must have defined property 'email'")
		}

		if err = checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}

		if user.FirstName == "" {
			return errors.New("User must have defined property 'first_name'")
		}

		if user.LastName == "" {
			return errors.New("User must have defined property 'last_name'")
		}

		return nil

	case "login":
		if user.Email == "" {
			return errors.New("User must have defined property 'email'")
		}

		if err = checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}

		if user.Password == "" {
			return errors.New("User must have defined property 'password'")
		}

		return nil

	default:
		if user.Username == "" {
			return errors.New("User must have defined property 'username'")
		}

		if user.Email == "" {
			return errors.New("User must have defined property 'email'")
		}

		if err = checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}

		if user.FirstName == "" {
			return errors.New("User must have defined property 'first_name'")
		}

		if user.LastName == "" {
			return errors.New("User must have defined property 'last_name'")
		}

		if user.Password == "" {
			return errors.New("User must have defined property 'password'")
		}

		return nil

	}
}
