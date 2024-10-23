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

type CreateLinkRequest struct {
	Title string `json:"title" binding:"required" example:"My GitHub"`
	URL   string `json:"url" binding:"required" example:"https://github.com/johndoe"`
}

func NewLinkHandler(linkService *services.LinkService) *LinkHandler {
	return &LinkHandler{LinkService: linkService}
}

// CreateLinkHandler godoc
// @Summary Create a new link
// @Description Create a new link for the authenticated user's profile
// @Tags links
// @Accept json
// @Produce json
// @Param link body CreateLinkRequest true "Link details"
// @Security BearerAuth
// @Success 201 "message: Link created successfully"
// @Failure 400 "error: Invalid input"
// @Failure 401 "error: Unauthorized"
// @Failure 400 "error: Link already exists"
// @Router /links [post]
func (h *LinkHandler) CreateLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var requestBody CreateLinkRequest
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

// UpdateLinkHandler godoc
// @Summary Update a link
// @Description Update an existing link for the authenticated user
// @Tags links
// @Accept json
// @Produce json
// @Param id path int true "Link ID" example(1)
// @Param link body CreateLinkRequest true "Updated link details"
// @Security BearerAuth
// @Success 200 "message: Link updated successfully"
// @Failure 400 "error: Invalid input"
// @Failure 401 "error: Unauthorized"
// @Failure 404 "error: Link not found"
// @Router /links/{id} [put]
func (h *LinkHandler) UpdateLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	linkId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	var updatedLink models.Link
	if err := c.ShouldBindJSON(&updatedLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.LinkService.UpdateLink(username.(string), linkId, updatedLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}

// DeleteLinkHandler godoc
// @Summary Delete a link
// @Description Delete an existing link from the authenticated user's profile
// @Tags links
// @Accept json
// @Produce json
// @Param id path int true "Link ID" example(1)
// @Security BearerAuth
// @Success 200 "message: Link deleted successfully"
// @Failure 400 "error: Invalid link ID"
// @Failure 401 "error: Unauthorized"
// @Failure 404 "error: Link not found"
// @Router /links/{id} [delete]
func (h *LinkHandler) DeleteLinkHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	linkId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	err = h.LinkService.DeleteLink(username.(string), linkId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}
