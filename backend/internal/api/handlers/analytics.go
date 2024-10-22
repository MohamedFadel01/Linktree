package handlers

import (
	"linktree-mohamedfadel-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	AnalyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{AnalyticsService: analyticsService}
}

// TrackLinkClickHandler godoc
// @Summary Track a link click
// @Description Records a click event for a specific link. If the request includes authentication, the click will be associated with the authenticated user.
// @Tags analytics
// @Accept json
// @Produce json
// @Param id path int true "Link ID" example(1)
// @Security BearerAuth
// @Success 200 "message: Click tracked successfully"
// @Failure 400 "error: Invalid link ID"
// @Failure 404 "error: Link not found"
// @Failure 500 "error: Internal server error"
// @Router /analytics/{id}/click [post]
func (h *AnalyticsHandler) TrackLinkClickHandler(c *gin.Context) {
	id := c.Param("id")
	linkId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid link ID"})
		return
	}

	var username string
	if usernameInterface, exists := c.Get("username"); exists {
		if usernameStr, ok := usernameInterface.(string); ok {
			username = usernameStr
		}
	}

	if err := h.AnalyticsService.TrackLinkClicks(uint64(linkId), username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Click tracked successfully"})
}
