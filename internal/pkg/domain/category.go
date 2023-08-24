package domain

type Category struct {
	MediaType string `json:"media_type,omitempty" bson:"media_type,omitempty"`
	Path      string `json:"path,omitempty" bson:"path,omitempty"`
	FileName  string `json:"file_name,omitempty" bson:"file_name,omitempty"`
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
