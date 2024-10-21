package services

import (
	"errors"
	"fmt"
	"linktree-mohamedfadel-backend/internal/models"
	"net/url"

	"gorm.io/gorm"
)

type LinkService struct {
	db *gorm.DB
}

func NewLinkService(db *gorm.DB) *LinkService {
	return &LinkService{db: db}
}

func (s *LinkService) CreateLink(username string, link models.Link) error {
	if link.Title == "" || link.URL == "" {
		return errors.New("required fields are missing")
	}

	if _, err := url.ParseRequestURI(link.URL); err != nil {
		return fmt.Errorf("invalid url")
	}

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	var existingLink models.Link
	if err := s.db.Where("url = ?", link.URL).First(&existingLink).Error; err == nil {
		return fmt.Errorf("link already exists")
	}

	newLink := models.Link{
		Title:  link.Title,
		URL:    link.URL,
		UserId: user.ID,
	}

	return s.db.Create(&newLink).Error
}

func (s *LinkService) UpdateLink(username string, linkId uint64, updatedLink models.Link) error {
	var link models.Link
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := s.db.Where("id = ? AND user_id = ?", linkId, user.ID).First(&link).Error; err != nil {
		return fmt.Errorf("link not found: %v", err)
	}

	if updatedLink.Title != "" {
		link.Title = updatedLink.Title
	}

	if updatedLink.URL != "" {
		if _, err := url.ParseRequestURI(updatedLink.URL); err != nil {
			return fmt.Errorf("invalid url")
		}
		link.URL = updatedLink.URL
	}

	return s.db.Save(&link).Error
}
