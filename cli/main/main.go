package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
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

func main() {

	if os.Args[1:][0] == "-if" || os.Args[1:][0] == "import-files" {
		importDailyIdFiles()
		return
	}

	var category, filePrefix, dbCollection, apiEndpoint string
	category = os.Args[1:][0]

	switch category {
	case "movie", "movies":
		category = "Movie"
		filePrefix = "movie"
		dbCollection = "movies"
		apiEndpoint = "movie"

	case "tvseries", "tv_series", "tv":
		category = "TVSeries"
		filePrefix = "tv_series"
		dbCollection = "tvseries"
		apiEndpoint = "tv"

	case "person", "people":
		category = "Person"
		filePrefix = "person"
		dbCollection = "people"
		apiEndpoint = "person"

	case "tvnetwork", "tvnetworks", "tv_network", "tv_networks":
		category = "TVNetwork"
		filePrefix = "tv_network"
		dbCollection = "tvnetworks"
		apiEndpoint = "network"

	case "collection", "collections":
		category = "Collection"
		filePrefix = "collection"
		dbCollection = "collections"
		apiEndpoint = "collection"

	case "keyword", "keywords":
		category = "Keyword"
		filePrefix = "keyword"
		dbCollection = "keywords"
		apiEndpoint = "keyword"

	default:
		log.Print("invalid option")
		return
	}

	startAt := 1
	limit := -1

	if len(os.Args[1:]) >= 2 {
		startAt, _ = strconv.Atoi(os.Args[1:][1])
	}

	if len(os.Args[1:]) >= 3 {
		limit, _ = strconv.Atoi(os.Args[1:][2])
	}

	today := time.Now().Format("01_02_2006")
	fileName := fmt.Sprintf("./daily_id_exports/%s_ids_%s.json.gz", filePrefix, today)

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

		if limit > 0 && count >= (startAt+limit) {
			break
		}

		wg.Add(1)

		go func(line string) {

			defer wg.Done()

			if count >= startAt {

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

				query := fmt.Sprintf("%s/%v", apiEndpoint, entity.ID)

				err = tmdb.Get(query, &entity.Data)
				if err != nil {
					logging.Error(err.Error())
					return
				}

				result, err := database.Upsert(entity, dbCollection)
				if err != nil {
					logging.Error(err.Error())
					return
				}

				logging.Info(fmt.Sprintf("%s %v - %v", strings.ToUpper(category), entity.ID, result))

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

	logging.Close()

}

func importDailyIdFiles() error {

	var err error

	for _, cat := range domain.CategoryList {
		logging.Info(fmt.Sprintf("%s - %s - %s", cat.MediaType, cat.FileName, file.GetFileName(cat.FileName)))

		fileName := file.GetFileName(cat.FileName)
		err := tmdb.FetchFileFromURL(fileName)

		if err != nil {
			return err
		}
	}

	return err
}
