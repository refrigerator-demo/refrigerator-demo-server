package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"unique_index;not null"`
	Email     string    `gorm:"unique_index;not null"`
	Password  string    `gorm:"not null"`
	Salted    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (user User) Validate() error {
	return nil
}

func (user *User) MakeHashedPassword() error {
	if 0 == len(user.Password) {
		return errors.New("Password must not be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if nil != err {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *User) CheckPasswordHashed(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err
}
