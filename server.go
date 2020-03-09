package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// BEGIN Structs
type User struct {
	UserName      string
	Discriminator string
	Token         string
	ImageURL      string
	Guilds        []UserGuild
}

type UserGuild struct {
	Name     string
	ID       string
	Enabled  bool
	Prefix   string
	Channels []Channel
}

// ENDOF Structs

// BEGIN Handlers
var user *User

/* binExecute receives a set of parameters from the http request following the standard URL query
bin?key=value. It then parses the key and if it corresponds to any known key, will execute and return to
index page with a status of 302 (See https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)
*/
func binExecute(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if request.Form["command"] != nil {
		if request.Form["command"][0] == "refresh" {
			refresh(client)
			updateConfigYaml()
			fmt.Println("Sucessfully refreshed")
			http.Redirect(writer, request, "/settings", 302)
		} else if request.Form["command"][0] == "changeGuildState" && request.Form["id"] != nil {
			id := request.Form["id"][0]
			var changedGuildSlice []Guild
			for _, guild := range conf.Session.Guilds {
				if guild.ID == id {
					guild.Enabled = !guild.Enabled
					changedGuildSlice = append(changedGuildSlice, guild)
					fmt.Println("Server with ID: " + id + "'s state has been changed sucessfully to: " + strconv.FormatBool(guild.Enabled))
					http.Redirect(writer, request, "/settings", 302)
				} else {
					changedGuildSlice = append(changedGuildSlice, guild)
				}
			}
			conf.Session.Guilds = changedGuildSlice
			updateConfigYaml()
		}
	}
}

func (user *User) updateGuilds() {
	var userGuilds []UserGuild
	for _, guild := range conf.Session.Guilds {
		var appendGuild UserGuild
		appendGuild.Name = guild.Name
		appendGuild.Enabled = guild.Enabled
		appendGuild.ID = guild.ID
		appendGuild.Prefix = guild.Prefix
		appendGuild.Channels = guild.Channels
		userGuilds = append(userGuilds, appendGuild)
	}
	user.Guilds = userGuilds
}

func (user *User) settingsHandler(writer http.ResponseWriter, request *http.Request) {
	user.updateGuilds()
	baseTemplate, _ := template.ParseFiles("./static/index.html")
	baseTemplate.Execute(writer, user)
}

func errorHandler(writer http.ResponseWriter, request *http.Request, status int) {
	writer.WriteHeader(status)
	fmt.Println("404")
	if status == http.StatusNotFound {
		fmt.Fprint(writer, "test")
	}
}

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	userInfo, _ := client.User("@me")
	var userGuilds []UserGuild
	for _, guild := range conf.Session.Guilds {
		var appendGuild UserGuild
		appendGuild.Name = guild.Name
		appendGuild.Enabled = guild.Enabled
		appendGuild.ID = guild.ID
		appendGuild.Prefix = guild.Prefix
		appendGuild.Channels = guild.Channels
		userGuilds = append(userGuilds, appendGuild)
	}
	user = &User{UserName: userInfo.Username, Discriminator: userInfo.Discriminator, Token: conf.Token, ImageURL: userInfo.AvatarURL(""), Guilds: userGuilds}
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/settings", user.settingsHandler)
	http.HandleFunc("/bin", binExecute)
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
