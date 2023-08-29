package domain

type Entity struct {
	ID   int         `json:"id,omitempty" bson:"id,omitempty"`
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
