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

func (s *LinkService) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return user, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

func (s *LinkService) CreateLink(link models.Link) error {
	if link.Title == "" || link.URL == "" {
		return errors.New("required fields are missing")
	}

	if _, err := url.ParseRequestURI(link.URL); err != nil {
		return fmt.Errorf("invalid url")
	}

	var existingLink models.Link
	if err := s.db.Where("url = ?", link.URL).First(&existingLink).Error; err == nil {
		return fmt.Errorf("link already exists")
	}

	return s.db.Create(&link).Error
}
