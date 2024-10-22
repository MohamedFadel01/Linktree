package models

import (
	"time"

	"gorm.io/datatypes"
)

// @Description Analytics data for tracking link usage
type Analytics struct {
	// ID is the unique identifier
	ID uint `gorm:"primarykey" json:"id" example:"1"`

	// ClickCount tracks number of clicks
	ClickCount uint `json:"click_count" example:"42"`

	// VisitorsUsernames stores usernames of visitors
	// swagger:strfmt json
	VisitorsUsernames datatypes.JSON `json:"visitors_usernames" swaggertype:"string" example:"[\"user1\", \"user2\"]"`

	// LinkID is the foreign key to the associated link
	LinkID uint `json:"link_id" example:"1"`

	// CreatedAt timestamp
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`

	// UpdatedAt timestamp
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}
