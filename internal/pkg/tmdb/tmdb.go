package tmdb

import (
	"encoding/json"
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/file"
	"gomdb/internal/pkg/logging"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type tmdbClient struct {
	_apiUrl string
	_apiKey string
}

func NewTmdbClient(apiUrl, apiKey string) (domain.EntityAPI, error) {

	client := tmdbClient{
		_apiUrl: apiUrl,
		_apiKey: apiKey,
	}

	err := client.PingAPI()

	return client, err
}

func (tc *tmdbClient) PingAPI() error {
	url := fmt.Sprintf("%s/configuration%s", tc._apiUrl, tc._apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Error(fmt.Sprintf("API ping failed: %s", resp.Status))
	}

	return nil
}

func (tc tmdbClient) GetFromAPI(entity *interface{}) error {
	entityType := reflect.TypeOf(*entity)
	entityTypeName := strings.ToLower(strings.Split(entityType.String(), ".")[1])
	id := reflect.ValueOf(*entity).Elem().FieldByName("ID").Int()
	err := tc.httpGet(fmt.Sprintf("%s/%v", entityTypeName, id), entity)
	return err
}

func FetchFileFromURL(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file.CopyFile(filePath, resp.Body)

	return err
}

func (tc tmdbClient) httpGet(query string, entity *interface{}) error {

	url := tc._apiUrl + query + tc._apiKey

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
