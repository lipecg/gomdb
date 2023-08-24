package domain

type ProductionCountry struct {
	Iso31661 string `json:"iso_3166_1,omitempty" bson:"iso_3166_1,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name,omitempty" bson:"english_name,omitempty"`
	Iso6391     string `json:"iso_639_1,omitempty" bson:"iso_639_1,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
}
