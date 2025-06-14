package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/protocol-golang/config"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/database"
	"github.com/sorrawichYooboon/protocol-golang/logger"
	"github.com/sorrawichYooboon/protocol-golang/migrations"
)

func main() {
	router := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database.Connect(cfg)
	migrations.RunMigrations(cfg)
	logger.InitLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	router.Run(":" + port)
}
