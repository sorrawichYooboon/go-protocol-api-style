package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/protocol-golang/config"
	"github.com/sorrawichYooboon/protocol-golang/graph"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/database"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/graphql"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/http"
	httphandler "github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/http/handler"
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

	router := gin.Default()

	movieHandler := httphandler.NewMovieHandler(movieUsecase)
	http.SetupRoutes(router, movieHandler)

	gqlResolver := &graph.Resolver{MovieUsecase: movieUsecase}
	router.POST("/graphql", graphql.GraphqlHandler(gqlResolver))
	router.GET("/playground", graphql.PlaygroundHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server running at :%s (REST, GraphQL, Playground)", port)
	router.Run(":" + port)
}
