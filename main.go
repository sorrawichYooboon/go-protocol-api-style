package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/protocol-golang/config"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/database"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/http"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/http/handler"
	"github.com/sorrawichYooboon/protocol-golang/internal/usecase"
	"github.com/sorrawichYooboon/protocol-golang/logger"
	"github.com/sorrawichYooboon/protocol-golang/migrations"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.Connect(cfg)
	migrations.RunMigrations(cfg)
	logger.InitLogger()

	movieRepo := database.NewMovieRepository(db)
	movieUsecase := usecase.NewMovieUsecase(movieRepo)
	movieHandler := handler.NewMovieHandler(movieUsecase)

	router := gin.Default()
	http.SetupRoutes(router, movieHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	router.Run(":" + port)
}
