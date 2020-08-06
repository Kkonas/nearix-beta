package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Permissions struct {
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

type Auth struct {
	Permissions *Permissions
	Logins      int
}

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

func Login(name string, hash string) *Response {
	// This function should try to login to the Nearix server and return True if login was successful.
	fmt.Println("[DEBUG] Attempting to login...")
	apiUrl := "http://localhost:8080/auth"
	response, err := http.PostForm(apiUrl, url.Values{
		"name": {name},
		"hash": {hash},
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println(string(body))
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	responseStruct := new(Response)
	decoder.Decode(responseStruct)
	return responseStruct
}

func DecodeMagic(response *Response) *Permissions {
	permissions := new(Permissions)
	base64decoded, _ := base64.StdEncoding.DecodeString(response.Magic)
	realMagic := strings.Split(string(base64decoded), "9")
	if realMagic[0] == "" {
		return nil
	}
	magic, _ := strconv.Atoi(realMagic[0])
	ac := tierNameTovalue("auto_catcher")
	as := tierNameTovalue("auto_spammer")
	ccl := tierNameTovalue("custom_catch_list")
	dp := tierNameTovalue("discard_pokemons")
	at := tierNameTovalue("auto_trade")

	if magic&ac == ac {
		permissions.AutoCatcher = true
	}
	if magic&as == as {
		permissions.AutoSpammer = true
	}
	if magic&ccl == ccl {
		permissions.CustomCatchList = true
	}
	if magic&dp == dp {
		permissions.DiscardPokemons = true
	}
	if magic&at == at {
		permissions.AutoTrade = true
	}

	return permissions
}

func GetAuth(name string, hash string) *Auth {
	resp := Login(name, hash)
	auth := Auth{Permissions: DecodeMagic(resp), Logins: resp.Logins}
	return &auth
}
