package main

import (
	"encoding/json"
	"fmt"
	"gomdb/internal/pkg/database"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	////// ROUTES //////

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		var field string
		var value int
		if len(query) > 0 {
			field = query.Get("field")
			value, _ = strconv.Atoi(query.Get("value"))
		}

		var list []domain.Entity

		var collection = "movies"
		err := database.List(field, value, collection, &list)
		if err != nil {
			logging.Error(err.Error())
			http.Error(w, "Error querying DB", http.StatusInternalServerError)
			return
		}
		jsonBytes, err := json.Marshal(list)

		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	})

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {

		var entity domain.Entity
		var err error
		entity.ID, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			fmt.Fprintf(w, "Invalid ID")
			return
		}
		var collection = "movies"
		database.Get(&entity, collection)

		jsonBytes, err := json.Marshal(entity)
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)

	})

	r.HandleFunc("/tvseries", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TV SHOWS weeee!")
	})

	r.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Peoplesssss!")
	})

	r.HandleFunc("/tvnetworks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TV Networks!")
	})

	r.HandleFunc("/keywords", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Keywords!")
	})

	fmt.Println("Listening on port 8181...")
	http.ListenAndServe("localhost:8181", r)

}
