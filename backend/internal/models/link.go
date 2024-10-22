package models

import "time"

// @Description A link entry with associated analytics
type Link struct {
	// ID is the unique identifier
	ID uint `gorm:"primarykey" json:"id" example:"1"`

	// Title of the link
	Title string `json:"title" example:"My GitHub Profile"`

	// URL of the link
	URL string `json:"url" gorm:"unique" example:"https://github.com/username"`

	// UserID is the foreign key to the owner
	UserID uint `json:"user_id" example:"1"`

	// Analytics data for this link
	Analytics Analytics `json:"analytics" gorm:"foreignKey:LinkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`

	// UpdatedAt timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
