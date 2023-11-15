package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
