package httphandler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichYooboon/protocol-golang/internal/dto"
	"github.com/sorrawichYooboon/protocol-golang/internal/usecase"
)

type MovieHandlerImpl struct {
	movieUsecase usecase.MovieUsecase
}

func NewMovieHandler(movieUsecase usecase.MovieUsecase) MovieHandler {
	return &MovieHandlerImpl{movieUsecase: movieUsecase}
}

func (h *MovieHandlerImpl) GetMovies(c *gin.Context) {
	movies, err := h.movieUsecase.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch movies"})
		return
	}

	resp := make([]dto.MovieResponse, 0, len(movies))
	for _, m := range movies {
		resp = append(resp, dto.MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			ReleaseDate: m.ReleaseDate,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MovieHandlerImpl) GetMovieByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	movie, err := h.movieUsecase.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch movie"})
		return
	}
	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	resp := dto.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		ReleaseDate: movie.ReleaseDate,
	}
	c.JSON(http.StatusOK, resp)
}
