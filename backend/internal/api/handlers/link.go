package handlers

import (
	"linktree-mohamedfadel-backend/internal/models"
	"linktree-mohamedfadel-backend/internal/services"
	"net/http"
	"strconv"

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

	var requestBody struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newLink := models.Link{
		Title: requestBody.Title,
		URL:   requestBody.URL,
	}

	if err := h.LinkService.CreateLink(username.(string), newLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Link created successfully"})
}

func (h *LinkHandler) UpdateLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	linkId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedLink models.Link
	if err := c.ShouldBindJSON(&updatedLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.LinkService.UpdateLink(username.(string), linkId, updatedLink)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}

func (h *LinkHandler) DeleteLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	linkId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.LinkService.DeleteLink(username.(string), linkId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}
