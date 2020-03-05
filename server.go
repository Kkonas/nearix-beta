package main

import (
	"log"
	"net/http"
	"text/template"
)

// BEGIN Structs
type User struct {
	UserName      string
	Discriminator string
}

// ENDOF Structs

// BEGIN Handlers
func (user *User) indexHandler(writer http.ResponseWriter, request *http.Request) {
	baseTemplate, _ := template.ParseFiles("./static/index.html")
	UserName := user.UserName
	baseTemplate.Execute(writer, UserName)
}

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	userInfo, _ := client.User("@me")
	user := &User{UserName: userInfo.Username, Discriminator: userInfo.Discriminator}
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/settings", user.indexHandler)
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
