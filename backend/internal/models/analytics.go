package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Analytics struct {
	gorm.Model
	ClickCount        uint           `json:"click_count"`
	VisitorsUsernames datatypes.JSON `json:"visitors_usernames"`
	LinkID            uint           `json:"link_id"`
}
