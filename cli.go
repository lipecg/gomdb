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

	"gomdb/cli/internal/pkg/database"
	"gomdb/cli/internal/pkg/file"
	"gomdb/cli/internal/pkg/http"
	"gomdb/cli/internal/pkg/models"
)

const downloadURL = "http://files.tmdb.org/p/exports/"

func main() {

	startAt, _ := strconv.Atoi(os.Args[1:][0])
	limit, _ := strconv.Atoi(os.Args[1:][1])

	f, err := os.OpenFile("./import.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}

	defer f.Close()

	// IMPORTS ALL DAILY ID FILES
	// importDailyIdFiles()

	//log.Print(http.GetMovieFromAPI(165))

	// Open the gzipped file
	file, err := os.Open("./daily_id_exports/movie_ids_08_18_2023.json.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
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

	// Iterate through each line and print it
	//TODO: REMOVE THIS LIMIT
	count := 1
	// var movies []interface{}
	for scanner.Scan() {

		wg.Add(1)

		go func(line string) {

			defer wg.Done()

			if count >= startAt {

				// Acquire a semaphore to limit the number of concurrent API calls
				semaphore <- struct{}{}

				// Wait for the rate limiter to allow the API call
				<-rateLimit

				var movie models.Movie
				err := json.Unmarshal([]byte(line), &movie)
				if err != nil {
					log.Fatal(err)
					f.WriteString(fmt.Sprintf("%v %s %v \n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", err))
				}

				movie = http.GetMovieFromAPI(movie.ID)

				result := database.UpdateMovieDB(&movie)

				f.WriteString(fmt.Sprintf("%v %s %v - %s - Result %v \n", time.Now().Format("2006-01-02 15:04:05"), "INFO", movie.ID, movie.OriginalTitle, result))

				//log.Print(movie.ID, movie.OriginalTitle)

				// Release the semaphore
				<-semaphore

				//log.Printf("%v - %s movie added to slice", movie.ID, movie.OriginalTitle)
				// if _, err :=  err != nil {
				// 	log.Print(err)
				//
				// }
			}

			count++

		}(scanner.Text())

		if count > startAt+limit {
			break
		}
	}

	wg.Wait()

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		f.WriteString(fmt.Sprintf("%v %s %v \n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", err))
	}

}

// func logError(err error) {
// 	f.WriteString(fmt.Sprintf("%v %s %v \n", time.Now().Format("2006-01-02 15:04:05"), "ERROR", err))
// }

func importDailyIdFiles() error {

	var err error

	for _, cat := range models.CategoryList {
		log.Printf("%s - %s - %s", cat.MediaType, cat.FileName, file.GetFileName(cat.FileName))

		fileName := file.GetFileName(cat.FileName)
		fileURL := downloadURL + fileName
		filePath := "./daily_id_exports/" + fileName
		err := http.FetchFileFromURL(fileURL, filePath)

		if err != nil {
			return err
		}
		// log.Printf("%s downloaded to %s", fileURL, filePath)
	}

	return err
}
