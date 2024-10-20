package handlers

import (
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	var requestBody struct {
		FullName string `json:"full_name"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user.FullName = requestBody.FullName
	user.Username = requestBody.Username
	user.Bio = requestBody.Bio

	if err := h.UserService.SignUp(user, requestBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
