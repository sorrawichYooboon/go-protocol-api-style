package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/protocol-golang/internal/infrastructure/http/handler"
)

func SetupRoutes(router *gin.Engine, movieHandler handler.MovieHandler) {
	movies := router.Group("/movies")
	{
		movies.GET("", movieHandler.GetMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
	}
}
