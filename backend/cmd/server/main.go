package main

import (
	"fmt"
	"linktree-mohamedfadel-backend/internal/api"
	"linktree-mohamedfadel-backend/internal/database"
	"linktree-mohamedfadel-backend/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

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

	router.SetupRoutes(engine)

	engine.Run(":8188")

}
