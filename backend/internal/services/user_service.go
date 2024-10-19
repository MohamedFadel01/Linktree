package services

import (
	"errors"
	"fmt"
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *UserService) SignUp(user models.User, password string) error {
	if user.FullName == "" || user.Username == "" || password == "" {
		return errors.New("required fields are missing")
	}

	var existingUser models.User
	if err := s.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("username already exists: %v", err)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	return s.db.Create(&user).Error
}

func (s *UserService) Login(username, password string) (string, error) {
	var user models.User

	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	if err := CheckPassword(user.PasswordHash, password); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}
