package app

import (
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
)

type tvSeriesSvc struct {
	DB domain.TVSeriesDB
}

func NewTVSeriesSvc(db domain.TVSeriesDB) domain.TVSeriesSvc {
	return tvSeriesSvc{
		DB: db,
	}
}

func (ms tvSeriesSvc) Get(id int) (*domain.TVSeries, error) {
	movie, err := ms.DB.Get(id)
	return (*movie).(*domain.TVSeries), err
}

func (ms tvSeriesSvc) List(query string) ([]*domain.TVSeries, error) {
	tvseries, err := ms.DB.List(query)
	tvseriesSlice := make([]*domain.TVSeries, len(tvseries))
	for i, v := range tvseries {
		tvseries, ok := (*v).(*domain.TVSeries)
		if !ok {
			logging.Error(fmt.Sprintf("Failed to convert element at index %d to *domain.TVSeries\n", i))
			continue
		}
		tvseriesSlice[i] = tvseries
	}
	return tvseriesSlice, err
}

func (ms tvSeriesSvc) Upsert(tvSeries *domain.TVSeries) error {
	var tvSeriesDB interface{} = tvSeries
	return ms.DB.Upsert(&tvSeriesDB)
}
