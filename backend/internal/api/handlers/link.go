package handlers

import (
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkHandler struct {
	LinkService *services.LinkService
}

func NewLinkHandler(linkService *services.LinkService) *LinkHandler {
	return &LinkHandler{LinkService: linkService}
}

func (h *LinkHandler) CreateLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.LinkService.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var requestBody struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	link := models.Link{
		Title:  requestBody.Title,
		URL:    requestBody.URL,
		UserId: user.ID,
	}

	if err := h.LinkService.CreateLink(link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Link created successfully"})
}
