package http

import (
	"encoding/json"
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/file"
	"io"
	"log"
	"net/http"
)

func FetchFileFromURL(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file.CopyFile(filePath, resp.Body)

	return err
}

func httpGet(query string) []byte {

	tmdbURL := "https://api.themoviedb.org/3/"
	tmdbApiKey := "?api_key=bdd0d7bc1bd4ee8f7c6b5fa9dc5611c1"
	url := tmdbURL + query + tmdbApiKey

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	return body

}

func extractJson(body []byte) domain.Movie {
	var movie domain.Movie
	if err := json.Unmarshal(body, &movie); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	return movie
}

func GetMovieFromAPI(id int) domain.Movie {
	obj := httpGet(fmt.Sprintf("movie/%v", id))
	movie := extractJson(obj)
	return movie
}
