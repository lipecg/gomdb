package domain

type Entity struct {
	ID       int         `json:"id,omitempty" bson:"id,omitempty"`
	ObjectId interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	ImdbID   string      `json:"imdb_id,omitempty" bson:"imdb_id,omitempty"`
}

type EntityDB interface {
	Get(id int) (*interface{}, error)
	List(search string) ([]*interface{}, error)
	Upsert(e *interface{}) error
}

type EntityAPI interface {
	GetFromAPI(e *interface{}) error
	// ListFromAPI(search string) ([]*interface{}, error)
	// ListChangedFromAPI(search string) ([]*interface{}, error)
}

type MovieTvIndex struct {
	ID            int    `json:"id,omitempty" bson:"id,omitempty"`
	OriginalTitle string `json:"original_title,omitempty" bson:"original_title,omitempty"`
}

type Category struct {
	MediaType string `json:"media_type,omitempty" bson:"media_type,omitempty"`
	Path      string `json:"path,omitempty" bson:"path,omitempty"`
	FileName  string `json:"file_name,omitempty" bson:"file_name,omitempty"`
}

type CreatedBy struct {
	ID          int    `json:"id,omitempty" bson:"id,omitempty"`
	CreditID    string `json:"credit_id,omitempty" bson:"credit_id,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Gender      int    `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfilePath string `json:"profile_path,omitempty" bson:"profile_path,omitempty"`
}

type Genre struct {
	ID   int    `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
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

type ProductionCountry struct {
	Iso31661 string `json:"iso_3166_1,omitempty" bson:"iso_3166_1,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
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

type SpokenLanguage struct {
	EnglishName string `json:"english_name,omitempty" bson:"english_name,omitempty"`
	Iso6391     string `json:"iso_639_1,omitempty" bson:"iso_639_1,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
}

var CategoryList = []Category{
	{
		MediaType: "Movie",
		FileName:  "movie_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "TV Series",
		FileName:  "tv_series_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "People",
		FileName:  "person_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "Collections",
		FileName:  "collection_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "TV Networks",
		FileName:  "tv_network_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "Keywords",
		FileName:  "keyword_ids_01_02_2006.json.gz",
	},
	{
		MediaType: "Production Companies",
		FileName:  "movie_ids_01_02_2006.json.gz",
	},
}
