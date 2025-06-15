package repository

import "github.com/sorrawichYooboon/protocol-golang/internal/domain"

type MovieRepository interface {
	GetAll() ([]domain.Movie, error)
	GetByID(id int64) (*domain.Movie, error)
}
