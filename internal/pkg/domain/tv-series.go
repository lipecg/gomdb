package domain

type TVSeries struct {
	Adult               bool                `json:"adult,omitempty" bson:"adult,omitempty"`
	BackdropPath        string              `json:"backdrop_path,omitempty" bson:"backdrop_path,omitempty"`
	CreatedBy           []CreatedBy         `json:"created_by,omitempty" bson:"created_by,omitempty"`
	EpisodeRunTime      []int               `json:"episode_run_time,omitempty" bson:"episode_run_time,omitempty"`
	FirstAirDate        string              `json:"first_air_date,omitempty" bson:"first_air_date,omitempty"`
	Genres              []Genre             `json:"genres,omitempty" bson:"genres,omitempty"`
	Homepage            string              `json:"homepage,omitempty" bson:"homepage,omitempty"`
	ID                  int                 `json:"id,omitempty" bson:"id,omitempty"`
	InProduction        bool                `json:"in_production,omitempty" bson:"in_production,omitempty"`
	Languages           []string            `json:"languages,omitempty" bson:"languages,omitempty"`
	LastAirDate         string              `json:"last_air_date,omitempty" bson:"last_air_date,omitempty"`
	LastEpisodeToAir    LastEpisodeToAir    `json:"last_episode_to_air,omitempty" bson:"last_episode_to_air,omitempty"`
	Name                string              `json:"name,omitempty" bson:"name,omitempty"`
	NextEpisodeToAir    interface{}         `json:"next_episode_to_air,omitempty" bson:"next_episode_to_air,omitempty"`
	Networks            []TvNetwork         `json:"networks,omitempty" bson:"networks,omitempty"`
	NumberOfEpisodes    int                 `json:"number_of_episodes,omitempty" bson:"number_of_episodes,omitempty"`
	NumberOfSeasons     int                 `json:"number_of_seasons,omitempty" bson:"number_of_seasons,omitempty"`
	OriginCountry       []string            `json:"origin_country,omitempty" bson:"origin_country,omitempty"`
	OriginalLanguage    string              `json:"original_language,omitempty" bson:"original_language,omitempty"`
	OriginalName        string              `json:"original_name,omitempty" bson:"original_name,omitempty"`
	Overview            string              `json:"overview,omitempty" bson:"overview,omitempty"`
	Popularity          float64             `json:"popularity,omitempty" bson:"popularity,omitempty"`
	PosterPath          string              `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	ProductionCompanies []ProductionCompany `json:"production_companies,omitempty" bson:"production_companies,omitempty"`
	ProductionCountries []ProductionCountry `json:"production_countries,omitempty" bson:"production_countries,omitempty"`
	Seasons             []Season            `json:"seasons,omitempty" bson:"seasons,omitempty"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages,omitempty" bson:"spoken_languages,omitempty"`
	Status              string              `json:"status,omitempty" bson:"status,omitempty"`
	Tagline             string              `json:"tagline,omitempty" bson:"tagline,omitempty"`
	Type                string              `json:"type,omitempty" bson:"type,omitempty"`
	VoteAverage         float64             `json:"vote_average,omitempty" bson:"vote_average,omitempty"`
	VoteCount           int                 `json:"vote_count,omitempty" bson:"vote_count,omitempty"`
}

type TVSeriesSvc interface {
	Get(id int) (*TVSeries, error)
	List(query string) ([]*TVSeries, error)
	Upsert(m *TVSeries) error
}

type TVSeriesDB interface {
	EntityDB
}

type Season struct {
	AirDate      string  `json:"air_date,omitempty" bson:"air_date,omitempty"`
	EpisodeCount int     `json:"episode_count,omitempty" bson:"episode_count,omitempty"`
	ID           int     `json:"id,omitempty" bson:"id,omitempty"`
	Name         string  `json:"name,omitempty" bson:"name,omitempty"`
	Overview     string  `json:"overview,omitempty" bson:"overview,omitempty"`
	PosterPath   string  `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	SeasonNumber int     `json:"season_number,omitempty" bson:"season_number,omitempty"`
	VoteAverage  float64 `json:"vote_average,omitempty" bson:"vote_average,omitempty"`
}

type LastEpisodeToAir struct {
	ID             int     `json:"id,omitempty" bson:"id,omitempty"`
	Name           string  `json:"name,omitempty" bson:"name,omitempty"`
	Overview       string  `json:"overview,omitempty" bson:"overview,omitempty"`
	VoteAverage    float64 `json:"vote_average,omitempty" bson:"vote_average,omitempty"`
	VoteCount      int     `json:"vote_count,omitempty" bson:"vote_count,omitempty"`
	AirDate        string  `json:"air_date,omitempty" bson:"air_date,omitempty"`
	EpisodeNumber  int     `json:"episode_number,omitempty" bson:"episode_number,omitempty"`
	EpisodeType    string  `json:"episode_type,omitempty" bson:"episode_type,omitempty"`
	ProductionCode string  `json:"production_code,omitempty" bson:"production_code,omitempty"`
	Runtime        int     `json:"runtime,omitempty" bson:"runtime,omitempty"`
	SeasonNumber   int     `json:"season_number,omitempty" bson:"season_number,omitempty"`
	ShowID         int     `json:"show_id,omitempty" bson:"show_id,omitempty"`
	StillPath      string  `json:"still_path,omitempty" bson:"still_path,omitempty"`
}
