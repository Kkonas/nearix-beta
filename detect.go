// This file will receive a Pokemon Image and will return the appropriate pokemon.
// This will get called from main.go

package main

import (
	//"fmt"
	"net/http"
	"os"
	"io"
)

// BEGIN function definition
func receive(url string){
	err := Download("images/template.jpg","url")

}
func Download(path string, url string) error{
	response, err := http.Get(url)
	logErr(err)
	var output *os.File

	if _, err := os.Stat(path); os.IsNotExist(err){
		logErr(err)
		output, err = os.Create(path)
		logErr(err)
	}else{
		err = os.Remove(path)
		logErr(err)
		output, err = os.Create(path)
		logErr(err)
	}
	// copy contents of response to output
	_, err = io.Copy(output, response.Body)
	defer response.Body.Close()
	defer output.Close()
	return err
}
func Compare() string{}

func Hash(path string) string{
	output, err := os.Open(path)
	logErr(err)

}

// ENDOF function definition
