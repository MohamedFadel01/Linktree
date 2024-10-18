package models

import "gorm.io/gorm"

type Analytics struct {
	gorm.Model
	ClickCount      uint   `json:"click_count"`
	VisitorUsername string `json:"visitor_username"`
	LinkID          uint   `json:"link_id"`
}
