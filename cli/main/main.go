package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gomdb/internal/pkg/database"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/file"
	"gomdb/internal/pkg/logging"
	"gomdb/internal/pkg/tmdb"
)

type categoryProperties struct {
	filePrefix   string
	dbCollection string
	apiEndpoint  string
}

type params struct {
	action       string
	category     string
	categoryInfo categoryProperties
	startAt      int
	limit        int
}

var execParams params = params{startAt: 1, limit: -1}
var categoryPropertiesMap = map[string]categoryProperties{
	"movies":      {"movie", "movies", "movie"},
	"tvseries":    {"tv_series", "tvseries", "tv"},
	"people":      {"person", "people", "person"},
	"tvnetworks":  {"tv_network", "tvnetworks", "network"},
	"collections": {"collection", "collections", "collection"},
	"keywords":    {"keyword", "keywords", "keyword"},
}

var validCategories = []string{"movies", "tvseries", "tvnetworks", "people", "collections", "keywords"}
var validActions = []string{"import", "sync", "download-files"}

func setExecParams(execArgs []string) error {

	if len(execArgs) < 1 {
		return fmt.Errorf("missing action: valid options are %v", validActions)
	}

	if !sliceContains(validActions, execArgs[0]) {
		return fmt.Errorf("invalid action: valid options are %v", validActions)
	}

	execParams.action = strings.ToLower(execArgs[0])

	if execParams.action == "download-files" {
		return nil
	}

	var ok bool
	execParams.category = strings.ToLower(execArgs[1])
	execParams.categoryInfo, ok = categoryPropertiesMap[execParams.category]
	if !ok {
		return fmt.Errorf("invalid category: valid options are %v", validCategories)
	}

	if execParams.action == "sync" {
		return nil
	}

	var err error

	if len(os.Args[1:]) > 2 && os.Args[1:][2] != "" {
		execParams.startAt, err = strconv.Atoi(os.Args[1:][2])
		if err != nil {
			return err
		}
	}

	if len(os.Args[1:]) > 3 && os.Args[1:][3] != "" {
		execParams.limit, err = strconv.Atoi(os.Args[1:][3])
		if err != nil {
			return err
		}
	}

	return nil
}

func sliceContains(collection []string, s string) bool {
	for _, item := range collection {
		if item == s {
			return true
		}
	}
	return false
}

func init() {
	if err := setExecParams(os.Args[1:]); err != nil {
		logging.Panic(err.Error())
	}
}

func main() {

	startTime := time.Now()
	logging.Infoln(fmt.Sprintf("Starting %s %s %s", execParams.action, execParams.category, startTime.Format("2006-01-02 15:04:05")))

	defer func() {
		endTime := time.Now()
		logging.Infoln(fmt.Sprintf("Finished %s %s %s (%v)", execParams.action, execParams.category, endTime.Format("2006-01-02 15:04:05"), endTime.Sub(startTime)))
		logging.Close()
	}()

	var err error

	switch execParams.action {
	case "download-files":
		err = importDailyIdFiles()
	case "sync":
		err = syncData()
	case "import":
		err = importData()
	}

	if err != nil {
		logging.Error(err.Error())
	}

}

func syncData() error {

	var updatedEntities *[]domain.Entity = new([]domain.Entity)

	err := tmdb.GetUpdatedEntities(execParams.categoryInfo.apiEndpoint, updatedEntities, tmdb.SearchOptions{Page: 1})
	if err != nil {
		logging.Panic(err.Error())
	}

	countEntities := len(*updatedEntities)
	countUpdatedEntities := 0

	logging.Infoln(fmt.Sprintf("Found %d updated %s", countEntities, execParams.category))

	var wg sync.WaitGroup

	// Create a channel to limit the number of concurrent API calls
	concurrency := 40
	semaphore := make(chan struct{}, concurrency)

	// Create a rate limiter to respect the rate limit of 40 calls per second
	rateLimit := time.Tick(time.Second / 40)

	for _, entity := range *updatedEntities {

		wg.Add(1)

		go func(ent domain.Entity) {

			defer wg.Done()

			// Acquire a semaphore to limit the number of concurrent API calls
			semaphore <- struct{}{}

			// Wait for the rate limiter to allow the API call
			<-rateLimit

			query := fmt.Sprintf("%s/%v", execParams.categoryInfo.apiEndpoint, ent.ID)

			err := tmdb.Get(query, &ent.Data)
			if err != nil {
				logging.Error(err.Error())
				return
			}

			result, err := database.Upsert(&ent, execParams.categoryInfo.dbCollection)
			if err != nil {
				logging.Error(err.Error())
				return
			}

			countUpdatedEntities++
			completion := fmt.Sprintf("%d/%d", countUpdatedEntities, countEntities)

			logging.Infoln(fmt.Sprintf("%s %s %v - %v", strings.ToUpper(execParams.category), completion, ent.ID, result))

			// Release the semaphore
			<-semaphore

		}(entity)

	}

	wg.Wait()

	return nil
}

func importData() error {
	today := time.Now().Format("01_02_2006")
	fileName := fmt.Sprintf("./daily_id_exports/%s_ids_%s.json.gz", execParams.categoryInfo.filePrefix, today)

	// Open the gzipped file
	file, err := os.Open(fileName)
	if err != nil {
		logging.Panic(err.Error())
	}
	defer file.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		logging.Panic(err.Error())
	}
	defer gzReader.Close()

	// Create a scanner to read the decompressed data line by line
	scanner := bufio.NewScanner(gzReader)

	var wg sync.WaitGroup

	// Create a channel to limit the number of concurrent API calls
	concurrency := 40
	semaphore := make(chan struct{}, concurrency)

	// Create a rate limiter to respect the rate limit of 10 calls per second
	rateLimit := time.Tick(time.Second / 10)

	count := 1

	for scanner.Scan() {

		if execParams.limit > 0 && count >= (execParams.startAt+execParams.limit) {
			break
		}

		wg.Add(1)

		go func(line string) {

			defer wg.Done()

			if count >= execParams.startAt {

				// Acquire a semaphore to limit the number of concurrent API calls
				semaphore <- struct{}{}

				// Wait for the rate limiter to allow the API call
				<-rateLimit

				var entity *domain.Entity

				err := json.Unmarshal([]byte(line), &entity)
				if err != nil {
					logging.Error(fmt.Sprintf("%s %v \n", "ERROR", err))
					return
				}

				query := fmt.Sprintf("%s/%v", execParams.categoryInfo.apiEndpoint, entity.ID)

				err = tmdb.Get(query, &entity.Data)
				if err != nil {
					logging.Error(err.Error())
					return
				}

				result, err := database.Upsert(entity, execParams.categoryInfo.dbCollection)
				if err != nil {
					logging.Error(err.Error())
					return
				}

				logging.Infoln(fmt.Sprintf("%s %v - %v", strings.ToUpper(execParams.category), entity.ID, result))

				// Release the semaphore
				<-semaphore

			}

		}(scanner.Text())

		count++
	}

	wg.Wait()

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		logging.Panic(err.Error())
	}

	return nil
}

func importDailyIdFiles() error {

	for k, v := range categoryPropertiesMap {

		fileName := file.GetFileName(v.filePrefix + "_ids_01_02_2006.json.gz")
		err := tmdb.FetchFileFromURL(fileName)

		if err != nil {
			return fmt.Errorf("%s %s - %s", k, fileName, err.Error())
		}

		logging.Infoln(fmt.Sprintf("%s %s - OK", k, fileName))
	}

	return nil
}
