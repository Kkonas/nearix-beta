package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// BEGIN Handlers

// indexHandler handles main page index
func indexHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	t, _ := template.ParseFiles("static/html/index.html")
	t.Execute(writer, nil)
	fmt.Println(request.Form)
}

// settingsHandler handles settings page
func settingsHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Fprintf(writer, "Hello settings!")
}

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/settings", settingsHandler)
	http.Handle("static", http.StripPrefix("static", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
