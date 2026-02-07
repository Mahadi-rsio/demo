package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Email     string `json:"email" gorm:"unique;not null"`
	Name      string `json:"name" gorm:"not null"`
	Password  string `json:"password"`
	Role      string `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {

	u.ID = uuid.New().String()

	password := u.Password

	if password == "" {
		fmt.Println("Password not found")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	u.Password = string(hash)

	return
}
