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

type Genre struct {
	ID   int    `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
