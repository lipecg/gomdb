package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"gomdb/app"
	"gomdb/internal/pkg/database"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/file"
	"gomdb/internal/pkg/http"
	"gomdb/internal/pkg/logging"
)

func main() {

	if os.Args[1:][0] == "import-files" {
		importDailyIdFiles()
		return
	}

	category := os.Args[1:][0]

	switch category {
	case "movie", "movies":
		category = "movie"
	case "tvseries", "tv_series":
		category = "tv_series"
		logging.Panic("Category not supported yet.")
	case "person", "people":
		category = "person"
		logging.Panic("Category not supported yet.")
	case "tvnetwork", "tvnetworks", "tv_network", "tv_networks":
		category = "tv_network"
		logging.Panic("Category not supported yet.")
	case "collection", "collections":
		category = "collection"
		logging.Panic("Category not supported yet.")
	case "keyword", "keywords":
		category = "keyword"
		logging.Panic("Category not supported yet.")
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
	fileName := fmt.Sprintf("./daily_id_exports/%s_ids_%s.json.gz", category, today)

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

	db, err := database.NewMongoStore()
	if err != nil {
		logging.Panic(err.Error())
	}
	svc := app.NewMovieSvc(db)

	var wg sync.WaitGroup

	// Create a channel to limit the number of concurrent API calls
	concurrency := 40
	semaphore := make(chan struct{}, concurrency)

	// Create a rate limiter to respect the rate limit of 10 calls per second
	rateLimit := time.Tick(time.Second / 10)

	count := 1
	// var movies []interface{}
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

				var movie domain.Movie
				err := json.Unmarshal([]byte(line), &movie)
				if err != nil {
					logging.Info(fmt.Sprintf("%s %v \n", "ERROR", err))
				}

				movie = http.GetMovieFromAPI(movie.ID)

				err = svc.Upsert(&movie)
				if err != nil {
					logging.Panic(err.Error())
				}

				logging.Info(fmt.Sprintf("%s %v - %s - %v \n", "MOVIE", movie.ID, movie.OriginalTitle, movie.ObjectId))

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

	downloadURL := "http://files.tmdb.org/p/exports/"

	var err error

	for _, cat := range domain.CategoryList {
		logging.Info(fmt.Sprintf("%s - %s - %s", cat.MediaType, cat.FileName, file.GetFileName(cat.FileName)))

		fileName := file.GetFileName(cat.FileName)
		fileURL := downloadURL + fileName
		filePath := "./daily_id_exports/" + fileName
		err := http.FetchFileFromURL(fileURL, filePath)

		if err != nil {
			return err
		}
	}

	return err
}
