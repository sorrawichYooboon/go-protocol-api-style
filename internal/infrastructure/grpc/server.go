package grpc

import (
	"context"

	"github.com/sorrawichYooboon/go-protocol-api-style/internal/infrastructure/grpc/moviepb"
	"github.com/sorrawichYooboon/go-protocol-api-style/internal/usecase"
)

type MovieServer struct {
	moviepb.UnimplementedMovieServiceServer
	MovieUsecase usecase.MovieUsecase
}

func NewMovieServer(usecase usecase.MovieUsecase) *MovieServer {
	return &MovieServer{MovieUsecase: usecase}
}

func (s *MovieServer) GetMovie(ctx context.Context, req *moviepb.GetMovieRequest) (*moviepb.GetMovieResponse, error) {
	movie, err := s.MovieUsecase.GetMovieByID(req.Id)
	if err != nil || movie == nil {
		return nil, err
	}
	return &moviepb.GetMovieResponse{
		Movie: &moviepb.Movie{
			Id:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			ReleaseDate: movie.ReleaseDate,
		},
	}, nil
}

func (s *MovieServer) ListMovies(ctx context.Context, req *moviepb.ListMoviesRequest) (*moviepb.ListMoviesResponse, error) {
	movies, err := s.MovieUsecase.GetAllMovies()
	if err != nil {
		return nil, err
	}
	resp := &moviepb.ListMoviesResponse{}
	for _, m := range movies {
		resp.Movies = append(resp.Movies, &moviepb.Movie{
			Id:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			ReleaseDate: m.ReleaseDate,
		})
	}
	return resp, nil
}
