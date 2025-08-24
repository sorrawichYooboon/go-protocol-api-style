package usecase

import (
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/domain"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/repository"
)

type MovieUsecaseImpl struct {
	movieRepo repository.MovieRepository
}

func NewMovieUsecase(repo repository.MovieRepository) MovieUsecase {
	return &MovieUsecaseImpl{movieRepo: repo}
}

func (u *MovieUsecaseImpl) GetAllMovies() ([]domain.Movie, error) {
	return u.movieRepo.GetAll()
}

func (u *MovieUsecaseImpl) GetMovieByID(id int64) (*domain.Movie, error) {
	movie, err := u.movieRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if movie == nil {
		return nil, nil
	}
	return movie, nil
}
