package models

import (
	"time"

	"gorm.io/datatypes"
)

type Analytics struct {
	ID                uint           `gorm:"primarykey"`
	ClickCount        uint           `json:"click_count"`
	VisitorsUsernames datatypes.JSON `json:"visitors_usernames"`
	LinkID            uint           `json:"link_id"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
