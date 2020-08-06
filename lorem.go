package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// RandomInt returns a random seeded integer from 0, max
func RandomInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

// Lorem returns a randomly sized string of Bacon Ipsum
func Lorem() string {
	bacon, err := http.Get("https://baconipsum.com/api/?type=meat-and-filler&paras=2&format=text")
	if err != nil {
		fmt.Println(err)
	}
	defer bacon.Body.Close()
	body, _ := ioutil.ReadAll(bacon.Body)
	words := strings.Split(string(body), " ")
	amount := RandomInt(6)
	var finalWords []string
	i := 0
	for i <= amount {
		finalWords = append(finalWords, words[i])
		i++
	}
	finalString := strings.Join(finalWords, " ")
	return finalString
}
