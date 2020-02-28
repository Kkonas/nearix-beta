package main

import (
	"log"
	"net/http"
)

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
