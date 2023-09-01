package domain

import "time"

type Entity struct {
	ID                  int                 `json:"id,omitempty" bson:"id,omitempty"`
	ObjectId            interface{}         `json:"_id,omitempty" bson:"_id,omitempty"`
	ImdbID              string              `json:"imdb_id,omitempty" bson:"imdb_id,omitempty"`
	Adult               *bool               `json:"adult,omitempty" bson:"adult,omitempty"`
	BackdropPath        string              `json:"backdrop_path,omitempty" bson:"backdrop_path,omitempty"`
	BelongsToCollection interface{}         `json:"belongs_to_collection,omitempty" bson:"belongs_to_collection,omitempty"`
	Budget              int                 `json:"budget,omitempty" bson:"budget,omitempty"`
	Genres              []Genre             `json:"genres,omitempty" bson:"genres,omitempty"`
	Homepage            string              `json:"homepage,omitempty" bson:"homepage,omitempty"`
	OriginalLanguage    string              `json:"original_language,omitempty" bson:"original_language,omitempty"`
	OriginalTitle       string              `json:"original_title,omitempty" bson:"original_title,omitempty"`
	Overview            string              `json:"overview,omitempty" bson:"overview,omitempty"`
	Popularity          float64             `json:"popularity,omitempty" bson:"popularity,omitempty"`
	PosterPath          string              `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	ProductionCompanies []ProductionCompany `json:"production_companies,omitempty" bson:"production_companies,omitempty"`
	ProductionCountries []ProductionCountry `json:"production_countries,omitempty" bson:"production_countries,omitempty"`
	ReleaseDate         string              `json:"release_date,omitempty" bson:"release_date,omitempty"`
	Revenue             int                 `json:"revenue,omitempty" bson:"revenue,omitempty"`
	Runtime             int                 `json:"runtime,omitempty" bson:"runtime,omitempty"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages,omitempty" bson:"spoken_languages,omitempty"`
	Status              string              `json:"status,omitempty" bson:"status,omitempty"`
	Tagline             string              `json:"tagline,omitempty" bson:"tagline,omitempty"`
	Title               string              `json:"title,omitempty" bson:"title,omitempty"`
	Video               *bool               `json:"video,omitempty" bson:"video,omitempty"`
	VoteAverage         float64             `json:"vote_average,omitempty" bson:"vote_average,omitempty"`
	VoteCount           int                 `json:"vote_count,omitempty" bson:"vote_count,omitempty"`
	Updated             time.Time           `json:"updated,omitempty" bson:"updated,omitempty"`

	Data interface{} `json:"data,omitempty" bson:"data,omitempty"`
}

type EntityDB interface {
	Get(id int, collection string) (*interface{}, error)
	List(search string, collection string) ([]*interface{}, error)
	Upsert(e *interface{}, collection string) error
}

type EntityAPI interface {
	GetFromAPI(e *interface{}, endpoint string) error
	// ListFromAPI(search string) ([]*interface{}, error)
	// ListChangedFromAPI(search string) ([]*interface{}, error)
}

type MovieTvIndex struct {
	ID            int    `json:"id,omitempty" bson:"id,omitempty"`
	OriginalTitle string `json:"original_title,omitempty" bson:"original_title,omitempty"`
}

type Genre struct {
	ID   int    `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type SearchResult struct {
	Results      []Entity `json:"results"`
	Page         int      `json:"page"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}
