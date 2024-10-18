package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName string `json:"full_name"`
	Username string `json:"username" gorm:"unique"`
	Bio      string `json:"bio"`
	Links    []Link `json:"links" gorm:"foreignKey:UserId"`
}
