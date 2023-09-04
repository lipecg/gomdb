package main

import (
	"bufio"
	"compress/gzip"
	"context"
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

	"golang.org/x/time/rate"
)

type categoryProperties struct {
	filePrefix   string
	dbCollection string
	apiEndpoint  string
}

type params struct {
	action            string
	category          string
	categoryInfo      categoryProperties
	startAt           int
	limit             int
	requestsPerSecond int
	batchSize         int
}

var execParams params = params{startAt: 1, limit: -1, requestsPerSecond: -1, batchSize: 10}
var categoryPropertiesMap = map[string]categoryProperties{
	"movies":      {"movie", "movies", "movie"},
	"tvseries":    {"tv_series", "tvseries", "tv"},
	"people":      {"person", "people", "person"},
	"tvnetworks":  {"tv_network", "tvnetworks", "network"},
	"collections": {"collection", "collections", "collection"},
	"keywords":    {"keyword", "keywords", "keyword"},
}

var validCategories = []string{"movies", "tvseries", "tvnetworks", "people", "collections", "keywords"}
var validActions = []string{"import", "import-many", "sync", "download-files"}

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

	if len(os.Args[1:]) > 4 && os.Args[1:][4] != "" {
		execParams.requestsPerSecond, err = strconv.Atoi(os.Args[1:][4])
		if err != nil {
			return err
		}
	}

	if len(os.Args[1:]) > 5 && os.Args[1:][5] != "" {
		execParams.batchSize, err = strconv.Atoi(os.Args[1:][5])
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
	case "import-many":
		err = importMany()
	}

	if err != nil {
		logging.Error(err.Error())
	}

}

func syncData() error {

	var updatedEntities []domain.Entity
	var options tmdb.SearchOptions = tmdb.SearchOptions{Page: 1}

	remainingPages := 1

	logging.Infoln("Fetching updates...")

	for remainingPages > 0 {

		var err error

		remainingPages, err = tmdb.GetUpdatedEntities(execParams.categoryInfo.apiEndpoint, &updatedEntities, options)
		if err != nil {
			return err
		}
		if remainingPages > 0 {
			options.Page++
		}

		fmt.Print("\r", len(updatedEntities), " entities fetched")
	}

	countEntities := len(updatedEntities)

	fmt.Print("\r")
	logging.Infoln(fmt.Sprintf("%v entities changed since last sync", countEntities))

	var insertableEntities []interface{}

	logging.Infoln("Fetching updated entities...")

	var wg sync.WaitGroup
	if execParams.requestsPerSecond < 0 {
		execParams.requestsPerSecond = 45
	}
	limiter := rate.NewLimiter(rate.Every(time.Second/time.Duration(execParams.requestsPerSecond)), 1)

	countEntities = 0

	for _, v := range updatedEntities {

		wg.Add(1)

		go func(v domain.Entity) {
			defer wg.Done()

			query := fmt.Sprintf("%s/%v", execParams.categoryInfo.apiEndpoint, v.ID)

			limiter.Wait(context.Background())
			err := tmdb.Get(query, &v)
			if err != nil {
				fmt.Print("x")
				return
			}

			v.Updated = time.Now()

			insertableEntities = append(insertableEntities, v)

			countEntities++

			fmt.Printf("\r%v entities fetched", countEntities)

		}(v)
	}

	wg.Wait()
	fmt.Println()

	countEntities = len(insertableEntities)

	fmt.Print("\r")
	logging.Infoln(fmt.Sprintf("%v entities to insert", countEntities))

	logging.Infoln("Inserting updated entities...")

	result, err := database.InsertMany(insertableEntities, execParams.categoryInfo.dbCollection)
	if err != nil {
		return err
	}

	logging.Infoln(fmt.Sprintf("%s %v entities inserted.", strings.ToUpper(execParams.category), len(result.InsertedIDs)))

	return nil
}

func importMany() error {
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

	// Create a rate limiter that allows X events per second
	var limiter *rate.Limiter = nil
	if execParams.requestsPerSecond > -1 {
		limiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(execParams.requestsPerSecond)), 1)
	}

	count := 1
	countEntities := 0
	var entities []interface{}

	for scanner.Scan() {

		line := scanner.Text()

		if execParams.limit > 0 && count >= (execParams.startAt+execParams.limit) {
			break
		}

		if count >= execParams.startAt {

			var entity *domain.Entity

			err := json.Unmarshal([]byte(line), &entity)
			if err != nil {
				logging.Error(fmt.Sprintf("%s %v \n", "ERROR", err))
				continue
			}

			entity.Updated = time.Now()

			entities = append(entities, entity)

			countEntities++

			if countEntities >= execParams.batchSize || (execParams.limit > 0 && count == (execParams.startAt+execParams.limit-1)) {

				wg.Add(1)

				go func(entitiesInsert []interface{}) {

					defer wg.Done()

					if limiter != nil {
						limiter.Wait(context.Background())
					}

					_, err := database.InsertMany(entitiesInsert, execParams.categoryInfo.dbCollection)
					if err != nil {
						logging.Error(err.Error())
						return
					}

					logging.Infoln(fmt.Sprintf("%s %v entities inserted.", strings.ToUpper(execParams.category), len(entitiesInsert)))

				}(entities)

				countEntities = 0
				entities = nil
			}

		}

		count++
	}

	wg.Wait()

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		logging.Panic(err.Error())
	}

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

	// Create a rate limiter that allows X events per second
	var limiter *rate.Limiter = nil
	if execParams.requestsPerSecond > -1 {
		limiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(execParams.requestsPerSecond)), 1)
	}

	count := 1

	for scanner.Scan() {

		if execParams.limit > 0 && count >= (execParams.startAt+execParams.limit) {
			break
		}

		wg.Add(1)

		go func(line string) {

			defer wg.Done()

			if count >= execParams.startAt {

				var entity *domain.Entity
				// {"adult":false,"id":9298,"original_title":"Ali G Indahouse","popularity":11.404,"video":false}

				err := json.Unmarshal([]byte(line), &entity)
				if err != nil {
					logging.Error(fmt.Sprintf("%s %v \n", "ERROR", err))
					return
				}
				// query := fmt.Sprintf("%s/%v", execParams.categoryInfo.apiEndpoint, entity.ID)

				if limiter != nil {
					limiter.Wait(context.Background())
				}

				// err = tmdb.Get(query, &entity.Data)
				// if err != nil {
				// 	logging.Error(err.Error())
				// 	return
				// }

				result, err := database.Upsert(entity, execParams.categoryInfo.dbCollection)
				if err != nil {
					logging.Error(err.Error())
					return
				}

				logging.Infoln(fmt.Sprintf("%s %v - %v", strings.ToUpper(execParams.category), entity.ID, result))
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
