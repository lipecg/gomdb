package tmdb

import (
	"encoding/json"
	"fmt"
	"gomdb/internal/pkg/file"
	"log"
	"net/http"
)

const tmdbURL = "https://api.themoviedb.org/3/"
const tmdbApiKey = "?api_key=bdd0d7bc1bd4ee8f7c6b5fa9dc5611c1"
const fileDownloadURL = "http://files.tmdb.org/p/exports/"
const fileDownloadDir = "./daily_id_exports/"

func Get(query string, entity *interface{}) error {

	url := tmdbURL + query + tmdbApiKey

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("search query failed: %s", resp.Status)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&entity)

	return err

}

func FetchFileFromURL(fileName string) error {
	resp, err := http.Get(fileDownloadURL + fileName)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	filePath := fileDownloadDir + fileName

	file.CopyFile(filePath, resp.Body)

	return err
}
