package models

type TvNetwork struct {
	Headquarters  string `json:"headquarters"`
	Homepage      string `json:"homepage"`
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}
