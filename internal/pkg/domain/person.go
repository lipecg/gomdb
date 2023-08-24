package domain

type Person struct {
	Adult              bool     `json:"adult,omitempty" bson:"adult,omitempty"`
	AlsoKnownAs        []string `json:"also_known_as,omitempty" bson:"also_known_as,omitempty"`
	Biography          string   `json:"biography,omitempty" bson:"biography,omitempty"`
	Birthday           string   `json:"birthday,omitempty" bson:"birthday,omitempty"`
	Deathday           string   `json:"deathday,omitempty" bson:"deathday,omitempty"`
	Gender             int      `json:"gender,omitempty" bson:"gender,omitempty"`
	Homepage           string   `json:"homepage,omitempty" bson:"homepage,omitempty"`
	ID                 int      `json:"id,omitempty" bson:"id,omitempty"`
	ImdbID             string   `json:"imdb_id,omitempty" bson:"imdb_id,omitempty"`
	KnownForDepartment string   `json:"known_for_department,omitempty" bson:"known_for_department,omitempty"`
	Name               string   `json:"name,omitempty" bson:"name,omitempty"`
	PlaceOfBirth       string   `json:"place_of_birth,omitempty" bson:"place_of_birth,omitempty"`
	Popularity         float64  `json:"popularity,omitempty" bson:"popularity,omitempty"`
	ProfilePath        string   `json:"profile_path,omitempty" bson:"profile_path,omitempty"`
}

type CreatedBy struct {
	ID          int    `json:"id,omitempty" bson:"id,omitempty"`
	CreditID    string `json:"credit_id,omitempty" bson:"credit_id,omitempty"`
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Gender      int    `json:"gender,omitempty" bson:"gender,omitempty"`
	ProfilePath string `json:"profile_path,omitempty" bson:"profile_path,omitempty"`
}
