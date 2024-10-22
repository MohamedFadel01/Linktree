package models

import "time"

type Link struct {
	ID        uint      `gorm:"primarykey"`
	Title     string    `json:"title"`
	URL       string    `json:"url" gorm:"unique"`
	UserID    uint      `json:"user_id"`
	Analytics Analytics `json:"analytics" gorm:"foreignKey:LinkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
