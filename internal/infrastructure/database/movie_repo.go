package database

import (
	"github.com/sorrawichYooboon/protocol-golang/internal/domain"
	"github.com/sorrawichYooboon/protocol-golang/internal/repository"
	"gorm.io/gorm"
)

type MovieRepositoryImpl struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) repository.MovieRepository {
	return &MovieRepositoryImpl{db: db}
}

func (r *MovieRepositoryImpl) GetAll() ([]domain.Movie, error) {
	var movies []domain.Movie
	if err := r.db.Table("movies").Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepositoryImpl) GetByID(id int64) (*domain.Movie, error) {
	var movie domain.Movie
	if err := r.db.Table("movies").Where("id = ?", id).First(&movie).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &movie, nil
}
