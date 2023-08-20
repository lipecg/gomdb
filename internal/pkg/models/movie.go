package models

import "time"

type Movie struct {
	ObjectId            interface{}         `json:"_id,omitempty" bson:"_id,omitempty"`
	Adult               bool                `json:"adult,omitempty" bson:"adult,omitempty"`
	BackdropPath        string              `json:"backdrop_path,omitempty" bson:"backdrop_path,omitempty"`
	BelongsToCollection interface{}         `json:"belongs_to_collection,omitempty" bson:"belongs_to_collection,omitempty"`
	Budget              int                 `json:"budget,omitempty" bson:"budget,omitempty"`
	Genres              []Genre             `json:"genres,omitempty" bson:"genres,omitempty"`
	Homepage            string              `json:"homepage,omitempty" bson:"homepage,omitempty"`
	ID                  int                 `json:"id,omitempty" bson:"id,omitempty"`
	ImdbID              string              `json:"imdb_id,omitempty" bson:"imdb_id,omitempty"`
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
	Video               bool                `json:"video,omitempty" bson:"video,omitempty"`
	VoteAverage         float64             `json:"vote_average,omitempty" bson:"vote_average,omitempty"`
	VoteCount           int                 `json:"vote_count,omitempty" bson:"vote_count,omitempty"`
	Updated             time.Time           `json:"updated,omitempty" bson:"updated,omitempty"`
}

// type Movie struct {
// 	ObjectId interface{} `bson:"_id,omitempty"`
// 	TmdbId   int         `json:"id,omitempty"`
// 	ImdbID   string      `json:"imdbID,omitempty" bson:"imdbID"`
// 	Title    string      `json:"Title,omitempty" bson:"title"`
// 	Year     string      `json:"Year,omitempty" bson:"year"`
// 	Rate     string      `json:"Rated,omitempty" bson:"rated"`
// 	Released string      `json:"Released,omitempty" bson:"released"`
// 	Runtime  string      `json:"Runtime,omitempty" bson:"runtime"`
// 	Genre    string      `json:"Genre,omitempty" bson:"genre"`
// 	Director string      `json:"Director,omitempty" bson:"director"`
// 	Writer   string      `json:"Writer,omitempty" bson:"writer"`
// 	Actors   string      `json:"Actors,omitempty" bson:"actors"`
// 	Plot     string      `json:"Plot,omitempty" bson:"plot"`
// 	Language string      `json:"Language,omitempty" bson:"language"`
// 	Country  string      `json:"Country,omitempty" bson:"country"`
// 	Awards   string      `json:"Awards,omitempty" bson:"awards"`
// 	Poster   string      `json:"Poster,omitempty" bson:"poster"`
// 	Ratings  []struct {
// 		Source string `json:"Source,omitempty" bson:"source"`
// 		Value  string `json:"Value,omitempty" bson:"value"`
// 	} `json:"Ratings,omitempty" bson:"ratings"`
// 	Metascore  string    `json:"Metascore,omitempty" bson:"metascore"`
// 	ImdbRating string    `json:"imdbRating,omitempty" bson:"imdbRating"`
// 	ImdbVotes  string    `json:"imdbVotes,omitempty" bson:"imdbVotes"`
// 	Type       string    `json:"Type,omitempty" bson:"type"`
// 	DVD        string    `json:"DVD,omitempty" bson:"dvd"`
// 	BoxOffice  string    `json:"BoxOffice,omitempty" bson:"boxOffice"`
// 	Production string    `json:"Production,omitempty" bson:"production"`
// 	Website    string    `json:"Website,omitempty" bson:"website"`
// 	Response   string    `json:"Response,omitempty" bson:"response"`
// 	Created    time.Time `json:"Created,omitempty" bson:"created"`
// 	Updated    time.Time `json:"Updated,omitempty" bson:"updated"`
// }
