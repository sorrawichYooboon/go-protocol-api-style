package usecase

import "github.com/sorrawichYooboon/go-protocol-api-style/internal/domain"

type MovieUsecase interface {
	GetAllMovies() ([]domain.Movie, error)
	GetMovieByID(id int64) (*domain.Movie, error)
}
