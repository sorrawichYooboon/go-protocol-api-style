package repository

import "github.com/sorrawichYooboon/go-protocol-api-style/internal/domain"

type MovieRepository interface {
	GetAll() ([]domain.Movie, error)
	GetByID(id int64) (*domain.Movie, error)
}
