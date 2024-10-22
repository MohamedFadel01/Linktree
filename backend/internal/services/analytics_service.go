package services

import (
	"encoding/json"
	"fmt"
	"linktree-mohamedfadel-backend/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) TrackLinkClicks(linkId uint64, visitorUsername string) error {
	var link models.Link
	if err := s.db.Where("id = ?", linkId).First(&link).Error; err != nil {
		return fmt.Errorf("link not found: %v", err)
	}

	var analytics models.Analytics
	if err := s.db.Where("link_id = ?", linkId).First(&analytics).Error; err != nil {
		emptyVisitors := make([]string, 0)
		visitorsJSON, err := json.Marshal(emptyVisitors)
		if err != nil {
			return fmt.Errorf("failed to marshal empty visitors array: %v", err)
		}

		analytics = models.Analytics{
			LinkID:            uint(linkId),
			ClickCount:        1,
			VisitorsUsernames: datatypes.JSON(visitorsJSON),
		}
	} else {
		analytics.ClickCount++
	}

	if visitorUsername != "" {
		var visitors []string
		if err := json.Unmarshal(analytics.VisitorsUsernames, &visitors); err != nil {
			return fmt.Errorf("failed to unmarshal visitors: %v", err)
		}

		userExists := false
		for _, visitor := range visitors {
			if visitor == visitorUsername {
				userExists = true
				break
			}
		}

		if !userExists {
			visitors = append(visitors, visitorUsername)
			visitorsJSON, err := json.Marshal(visitors)
			if err != nil {
				return fmt.Errorf("failed to marshal visitors: %v", err)
			}
			analytics.VisitorsUsernames = datatypes.JSON(visitorsJSON)
		}
	}

	if err := s.db.Save(&analytics).Error; err != nil {
		return fmt.Errorf("failed to save analytics: %v", err)
	}

	return nil
}
