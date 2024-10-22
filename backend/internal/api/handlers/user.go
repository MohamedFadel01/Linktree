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

type SignUpRequest struct {
	FullName string `json:"full_name" binding:"required" example:"John Doe"`
	Username string `json:"username" binding:"required" example:"johndoe"`
	Bio      string `json:"bio" example:"Software Developer | Tech Enthusiast"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// SignUpHandler godoc
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body SignUpRequest true "User registration details"
// @Success 201 "message: User created successfully"
// @Failure 400 "error: Invalid input"
// @Failure 400 "error: Username already exists"
// @Router /users/signup [post]
func (h *UserHandler) SignUpHandler(c *gin.Context) {
	var user models.User
	var requestBody SignUpRequest

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

// LoginHandler godoc
// @Summary Login user
// @Description Authenticate user credentials and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 "token: JWT_TOKEN_STRING"
// @Failure 400 "error: Invalid input"
// @Failure 401 "error: Invalid username or password"
// @Router /users/login [post]
func (h *UserHandler) LoginHandler(c *gin.Context) {
	var requestBody LoginRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.UserService.Login(requestBody.Username, requestBody.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetUserProfileInfoHandler godoc
// @Summary Get user profile
// @Description Retrieve user profile information and their associated links
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Username" example:"johndoe"
// @Success 200 {object} models.User "User profile with associated links"
// @Failure 404 "error: User not found"
// @Router /users/{username} [get]
func (h *UserHandler) GetUserProfileInfoHandler(c *gin.Context) {
	username := c.Param("username")

	user, err := h.UserService.GetUserProfileInfo(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserHandler godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Updated user information"
// @Security BearerAuth
// @Success 200 "message: User updated successfully"
// @Failure 400 "error: Invalid input"
// @Failure 401 "error: Unauthorized"
// @Failure 404 "error: User not found"
// @Router /users [put]
func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.UserService.UpdateUser(username.(string), updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUserHandler godoc
// @Summary Delete user account
// @Description Permanently delete the authenticated user's account and all associated data
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 "message: User deleted successfully"
// @Failure 401 "error:Unauthorized"
// @Failure 404 "error:User not found"
// @Failure 500 "error:Internal server error"
// @Router /users [delete]
func (h *UserHandler) DeleteUserHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.UserService.DeleteUser(username.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
