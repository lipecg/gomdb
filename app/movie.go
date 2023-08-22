package app

import "gomdb/internal/pkg/domain"

type movieSvc struct {
	DB domain.MovieDB
}

func NewMovieSvc(db domain.MovieDB) domain.MovieSvc {
	return movieSvc{
		DB: db,
	}
}

func (ms movieSvc) Get(id int) (*domain.Movie, error) {
	return ms.DB.Get(id)
}

func (ms movieSvc) List(query string) ([]*domain.Movie, error) {
	return ms.DB.List(query)
}

func (ms movieSvc) Upsert(movie *domain.Movie) error {
	return ms.DB.Upsert(movie)
}
