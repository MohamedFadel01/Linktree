package main

import (
	"fmt"
	"linktree-mohamedfadel-backend/internal/database"
	"log"
)

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to database successfully✅✅✅")
	}
}
