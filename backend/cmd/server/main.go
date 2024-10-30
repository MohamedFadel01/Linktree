package main

import (
	"fmt"
	"linktree-mohamedfadel-backend/docs"
	"linktree-mohamedfadel-backend/internal/api"
	"linktree-mohamedfadel-backend/internal/database"
	"linktree-mohamedfadel-backend/internal/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           Linktree API
// @version         1.0
// @description     A Linktree clone API server.

// @host      localhost:8188
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to database successfully✅✅✅")
	}

	userService := services.NewUserService(database.DB)
	linkService := services.NewLinkService(database.DB)
	analyticsService := services.NewAnalyticsService(database.DB)

	router := api.NewRouter(userService, linkService, analyticsService)

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"

	router.SetupRoutes(engine)

	engine.Run(":8188")
}
