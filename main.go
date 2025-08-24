package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/go-protocol-api-style/config"
	"github.com/sorrawichYooboon/go-protocol-api-style/graph"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/database"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/graphql"
	grpcinfra "github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/grpc"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/grpc/moviepb"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/http"
	httphandler "github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/http/handler"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/soap"
	soaphandler "github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/soap/handler"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/usecase"
	"github.com/sorrawichYooboon/go-protocol-api-style/logger"
	"github.com/sorrawichYooboon/go-protocol-api-style/migrations"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	movieSOAPHandler := soaphandler.NewMovieSOAPHandler(movieUsecase)
	soap.SetupSOAPRoutes(router, movieSOAPHandler)

	go func() {
		grpcPort := ":50051"
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen gRPC: %v", err)
		}
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer) // for development purposes, in production you might want to disable this because it exposes all services
		moviepb.RegisterMovieServiceServer(grpcServer, grpcinfra.NewMovieServer(movieUsecase))
		log.Printf("gRPC server running at %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server running at :%s (REST, GraphQL, Playground)", port)
	router.Run(":" + port)
}
