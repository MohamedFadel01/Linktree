package models

import "gorm.io/gorm"

type Link struct {
	gorm.Model
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	UserId    uint      `json:"user_id"`
	Analytics Analytics `json:"analytics" gorm:"foreignKey:LinkID;constrain:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
