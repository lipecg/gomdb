package main

import (
	"fmt"
	"net/http"
)

func main() {

	////// ROUTES //////

	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Movies! %s", r.Method)
	})

	http.HandleFunc("/tvseries", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TV SHOWS weeee!")
	})

	http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Peoplesssss!")
	})

	http.HandleFunc("/tvnetworks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TV Networks!")
	})

	http.HandleFunc("/keywords", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Keywords!")
	})

	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)

}
