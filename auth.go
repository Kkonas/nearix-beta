package main

import (
	"database/sql"
	"flag"
	"fmt"
	//"golang.org/x/crypto/bcrypt"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "nearixadmin"
	dbname = "nearix"
)

const (
	auto_catcher      = 0x1
	auto_spammer      = 0x2
	custom_catch_list = 0x4
	discard_pokemons  = 0x8
	auto_trade        = 0x10
)

func tierNameTovalue(name string) int {
	if name == "auto_catcher" {
		return auto_catcher
	} else if name == "auto_spammer" {
		return auto_spammer
	} else if name == "custom_catch_list" {
		return custom_catch_list
	} else if name == "discard_pokemons" {
		return discard_pokemons
	} else if name == "auto_trade" {
		return auto_trade
	}
	return 0
}

var sqlString string = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
var db *sql.DB
var err error

type User struct {
	Name          string
	Hash          string
	Version_name  string
	Register_date time.Time
	Logins        int
}

type Version struct {
	Name            string
	AutoCatcher     bool
	AutoSpammer     bool
	CustomCatchList bool
	DiscardPokemons bool
	AutoTrade       bool
}

type Response struct {
	Magic  string
	Logins int
}

func (user *User) String() string {
	return "[Name:" + space(user.Name) + ", Hash:" + space(user.Hash) + ", Version:" + space(user.Version_name) + ", Register:" + space(user.Register_date.String()) + ", Logins:" + space(strconv.Itoa(user.Logins)+"]")
}

func randMax(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func generateWithMaxSum(max int, length int) string {
	var sum int = 0
	var maxValue int = 8
	var returnRandom string
	for sum < max {
		random := randMax(maxValue)
		if random+sum <= max {
			sum = sum + random
			returnRandom = returnRandom + strconv.Itoa(random)
		} else {
			maxValue = max - sum
			continue
		}
	}
	if len(returnRandom) < length {
		returnRandom = returnRandom + "9"
		var count int = 0
		totalLength := len(returnRandom)
		for count <= (length - totalLength) {
			returnRandom = returnRandom + strconv.Itoa(randMax(10))
			count = count + 1
		}
	}
	returnRandom = base64.StdEncoding.EncodeToString([]byte(returnRandom))
	return returnRandom
}

func checkUser(name string, hash string) (*User, error) {
	user := new(User)
	row, err := db.Query(`SELECT * FROM nearix WHERE name = $1 AND hash = $2`, name, hash)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&user.Name, &user.Hash, &user.Version_name, &user.Register_date, &user.Logins)
		if err != nil {
			return nil, err
		}
		fmt.Println(user.String())
		return user, nil
	}
	return nil, err
}

func (user *User) addLogin() {
	user.Logins = user.Logins + 1
	_, err := db.Exec(`UPDATE nearix SET logins = $1 WHERE name = $2 AND hash = $3`, user.Logins, user.Name, user.Hash)
	if err != nil {
		return
	}
}

func space(st string) string {
	return " " + st
}

func calculatePermissionInteger(user *User) int {
	version := new(Version)
	versionRow, err := db.Query(`SELECT * FROM versions WHERE name=$1`, user.Version_name)
	var hexInteger int = 0
	if err != nil {
		log.Println(err)
		return 0
	}
	defer versionRow.Close()
	for versionRow.Next() {
		err = versionRow.Scan(&version.Name, &version.AutoCatcher, &version.AutoSpammer, &version.CustomCatchList, &version.DiscardPokemons, &version.AutoTrade)
		if err != nil {
			log.Println(err)
			return 0
		}
		if version.AutoCatcher {
			hexInteger = hexInteger + tierNameTovalue("auto_catcher")
		}
		if version.AutoSpammer {
			hexInteger = hexInteger + tierNameTovalue("auto_spammer")
		}
		if version.CustomCatchList {
			hexInteger = hexInteger + tierNameTovalue("custom_catch_list")
		}
		if version.DiscardPokemons {
			hexInteger = hexInteger + tierNameTovalue("discard_pokemons")
		}
		if version.AutoTrade {
			hexInteger = hexInteger + tierNameTovalue("auto_trade")
		}

	}
	return hexInteger
}

func authRouter(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	defer request.Body.Close()
	if request.Form["name"] != nil && request.Form["hash"] != nil {
		fmt.Println("Form:" + space(request.FormValue("name")) + space(request.FormValue("hash")))
		if err != nil {
			log.Println(err)
			return
		}
		// Check if user exists with given password hash
		foundUser, err := checkUser(request.FormValue("name"), request.FormValue("hash"))
		if err != nil || foundUser == nil {
			return
		}
		foundUser.addLogin()
		data := generateWithMaxSum(calculatePermissionInteger(foundUser), 30)
		response := Response{Magic: data, Logins: foundUser.Logins}
		json.NewEncoder(writer).Encode(response)

	}

}

func main() {
	port := flag.Int("port", 8080, "Select the auth server running port.")
	flag.Parse()
	db, err = sql.Open("postgres", sqlString)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Running auth server on port: " + strconv.Itoa(*port))
	http.HandleFunc("/auth", authRouter)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
