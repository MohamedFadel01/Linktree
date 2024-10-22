package api

import (
	"linktree-mohamedfadel-backend/internal/api/handlers"
	"linktree-mohamedfadel-backend/internal/api/middleware"
	"linktree-mohamedfadel-backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	userHandler      *handlers.UserHandler
	linkHandler      *handlers.LinkHandler
	analyticsHandler *handlers.AnalyticsHandler
}

func NewRouter(
	userService *services.UserService,
	linkService *services.LinkService,
	analyticsService *services.AnalyticsService,
) *Router {
	return &Router{
		userHandler:      handlers.NewUserHandler(userService),
		linkHandler:      handlers.NewLinkHandler(linkService),
		analyticsHandler: handlers.NewAnalyticsHandler(analyticsService),
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api/v1")
	{
		users := public.Group("/users")
		{
			users.POST("/signup", r.userHandler.SignUpHandler)
			users.POST("/login", r.userHandler.LoginHandler)
			users.GET("/:username", r.userHandler.GetUserProfileInfoHandler)
		}
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.ValidateJWTFromContext())
	{
		users := protected.Group("/users")
		{
			users.PUT("", r.userHandler.UpdateUserHandler)
			users.DELETE("", r.userHandler.DeleteUserHandler)
		}

		links := protected.Group("/links")
		{
			links.POST("", r.linkHandler.CreateLinkHandler)
			links.PUT("/:id", r.linkHandler.UpdateLinkHandler)
			links.DELETE("/:id", r.linkHandler.DeleteLinkHandler)
		}
	}

	optionalAuth := router.Group("/api/v1")
	optionalAuth.Use(middleware.OptionalJWTFromContext())
	{
		analytics := optionalAuth.Group("/analytics")
		{
			analytics.POST("/:id/click", r.analyticsHandler.TrackLinkClickHandler)
		}
	}
}
