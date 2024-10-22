package models

import "time"

// @Description A user account with profile information and associated links
type User struct {
	// ID is the unique identifier
	ID uint `gorm:"primarykey" json:"id" example:"1"`

	// FullName of the user
	FullName string `json:"full_name" example:"John Doe"`

	// Username is the unique identifier for the user
	Username string `json:"username" gorm:"unique" example:"johndoe"`

	// Bio contains user's description
	Bio string `json:"bio" example:"Software developer passionate about Go"`

	// Links associated with this user
	Links []Link `json:"links" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// PasswordHash stores the hashed password (not exposed in JSON)
	PasswordHash string `json:"-"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`

	// UpdatedAt timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
