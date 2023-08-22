package app

import (
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
)

type movieSvc struct {
	DB domain.MovieDB
}

func NewMovieSvc(db domain.MovieDB) domain.MovieSvc {
	return movieSvc{
		DB: db,
	}
}

func (ms movieSvc) Get(id int) (*domain.Movie, error) {
	movie, err := ms.DB.Get(id)
	return (*movie).(*domain.Movie), err
}

func (ms movieSvc) List(query string) ([]*domain.Movie, error) {
	movies, err := ms.DB.List(query)
	movieSlice := make([]*domain.Movie, len(movies))
	for i, v := range movies {
		movie, ok := (*v).(*domain.Movie)
		if !ok {
			logging.Error(fmt.Sprintf("Failed to convert element at index %d to *domain.Movie\n", i))
			continue
		}
		movieSlice[i] = movie
	}
	return movieSlice, err
}

func (ms movieSvc) Upsert(movie *domain.Movie) error {
	var movieDB interface{} = movie
	return ms.DB.Upsert(&movieDB)
}
