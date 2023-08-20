package models

type ProductionCompany struct {
	Description   string      `json:"description,omitempty" bson:"description,omitempty"`
	Headquarters  string      `json:"headquarters,omitempty" bson:"headquarters,omitempty"`
	Homepage      string      `json:"homepage,omitempty" bson:"homepage,omitempty"`
	ID            int         `json:"id,omitempty" bson:"id,omitempty"`
	LogoPath      string      `json:"logo_path,omitempty" bson:"logo_path,omitempty"`
	Name          string      `json:"name,omitempty" bson:"name,omitempty"`
	OriginCountry string      `json:"origin_country,omitempty" bson:"origin_country,omitempty"`
	ParentCompany interface{} `json:"parent_company,omitempty" bson:"parent_company,omitempty"`
}
