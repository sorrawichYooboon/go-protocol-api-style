package handler

import "github.com/gin-gonic/gin"

type MovieHandler interface {
	GetMovies(*gin.Context)
	GetMovieByID(*gin.Context)
}
