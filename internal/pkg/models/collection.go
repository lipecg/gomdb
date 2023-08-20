package models

type Collection struct {
	ID           int     `json:"id,omitempty" bson:"id,omitempty"`
	Name         string  `json:"name,omitempty" bson:"name,omitempty"`
	Overview     string  `json:"overview,omitempty" bson:"overview,omitempty"`
	PosterPath   string  `json:"poster_path,omitempty" bson:"poster_path,omitempty"`
	BackdropPath string  `json:"backdrop_path,omitempty" bson:"backdrop_path,omitempty"`
	Parts        []Movie `json:"parts,omitempty" bson:"parts,omitempty"`
}
