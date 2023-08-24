package app

import (
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
)

type movieSvc struct {
	DB  domain.MovieDB
	API domain.MovieAPI
}

func NewMovieSvc(db domain.MovieDB, api domain.MovieAPI) domain.MovieSvc {
	return movieSvc{
		DB:  db,
		API: api,
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

func (ms movieSvc) GetFromAPI(movie *domain.Movie) error {
	var movieAPI interface{} = movie
	return ms.API.GetFromAPI(&movieAPI)
}
