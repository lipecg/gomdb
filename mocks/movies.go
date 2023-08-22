package mocks

import "gomdb/internal/pkg/domain"

type MovieSvc struct {
	GetMovieResp   domain.Movie
	GetMovieErr    error
	ListMovieResp  []*domain.Movie
	ListMovieErr   error
	UpsertMovieErr error
}

func (ms MovieSvc) Get(id int) (*domain.Movie, error) {
	return &ms.GetMovieResp, ms.GetMovieErr
}

func (ms MovieSvc) List(query string) ([]*domain.Movie, error) {
	return ms.ListMovieResp, ms.ListMovieErr
}

func (ms MovieSvc) Upsert(*domain.Movie) error {
	return ms.UpsertMovieErr
}
