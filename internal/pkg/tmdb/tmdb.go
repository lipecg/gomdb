package tmdb

import (
	"encoding/json"
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/file"
	"net/http"
)

const tmdbURL = "https://api.themoviedb.org/3/"
const tmdbApiKey = "?api_key=bdd0d7bc1bd4ee8f7c6b5fa9dc5611c1"
const fileDownloadURL = "http://files.tmdb.org/p/exports/"
const fileDownloadDir = "./daily_id_exports/"

type SearchOptions struct {
	Page int
}

func Get(query string, entity *domain.Entity) error {

	url := tmdbURL + query + tmdbApiKey

	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("error http get: %s %s", url, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("search query failed: %s %s", url, resp.Status)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&entity)
	if err != nil {
		fmt.Print(resp.Body)
		return fmt.Errorf("error decoding object %v - %s", entity.ID, err.Error())
	}

	return nil

}

func GetUpdatedEntities(apiEndpoint string, updatedEntities *[]domain.Entity, options SearchOptions) (int, error) {

	remainingPages := 0

	url := tmdbURL + apiEndpoint + "/changes" + tmdbApiKey + "&page=" + fmt.Sprintf("%d", options.Page)

	resp, err := http.Get(url)

	if err != nil {
		return remainingPages, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return remainingPages, fmt.Errorf("search query failed: %s %s", url, resp.Status)
	}

	defer resp.Body.Close()

	var results domain.SearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return remainingPages, err
	}

	*updatedEntities = append(*updatedEntities, results.Results...)

	remainingPages = results.TotalPages - results.Page

	return remainingPages, err
}

func FetchFileFromURL(fileName string) error {
	resp, err := http.Get(fileDownloadURL + fileName)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search query failed: %s", resp.Status)
	}

	filePath := fileDownloadDir + fileName

	err = file.CopyFile(filePath, resp.Body)

	return err
}
