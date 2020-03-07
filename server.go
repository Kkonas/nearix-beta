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
	ImageURL      string
	Guilds        []UserGuild
}

type UserGuild struct {
	Name      string
	IsEnabled bool
}

// ENDOF Structs

// BEGIN Handlers
func (user *User) indexHandler(writer http.ResponseWriter, request *http.Request) {
	baseTemplate, _ := template.ParseFiles("./static/index.html")
	baseTemplate.Execute(writer, user)
}

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	userInfo, _ := client.User("@me")
	var userGuilds []UserGuild
	for _, guild := range client.State.Guilds {
		var appendGuild UserGuild
		appendGuild.Name = guild.Name
		appendGuild.IsEnabled = false
		userGuilds = append(userGuilds, appendGuild)
	}
	user := &User{UserName: userInfo.Username, Discriminator: userInfo.Discriminator, ImageURL: userInfo.AvatarURL(""), Guilds: userGuilds}
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/settings", user.indexHandler)
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
