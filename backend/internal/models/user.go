package models

import "time"

type User struct {
	ID           uint   `gorm:"primarykey"`
	FullName     string `json:"full_name"`
	Username     string `json:"username" gorm:"unique"`
	Bio          string `json:"bio"`
	Links        []Link `json:"links" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PasswordHash string `json:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
