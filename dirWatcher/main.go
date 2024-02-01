package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config := loadConfig()

	// Initialize database
	db := initDatabase()

	defer db.Close()

    log.Println("Database schema initialized successfully")

	// Initialize and start the watcher
	watcher := NewWatcher(config, db)
	go watcher.Start()

	// Initialize API
	router := gin.Default()
	api := NewAPI(config, db, watcher)
	api.RegisterRoutes(router)

	// Start API server
	apiPort := fmt.Sprintf(":%d", config.APIPort)
	err := router.Run(apiPort)
	if err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
