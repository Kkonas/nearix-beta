// This file will receive a Pokemon Image and will return the appropriate pokemon.
// This will get called from main.go

package main

import (
	//"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"io"
)

// BEGIN function definition
func receive(url string) string{
	err := Download("images/template.jpg",url)
	logErr(err)
	hash := Hash("images/template.jpg")
	return(hex.EncodeToString(hash)+"\n")
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
func Compare() string{return "Hello World!"}

func Hash(path string) []byte{
	output, err := ioutil.ReadFile(path)
	logErr(err)
	hash := md5.New()
	hash.Write(output)
	hashSum := hash.Sum(nil)
	hexEncoded := make([]byte, hex.EncodedLen(len(hashSum)))
	hex.Encode(hexEncoded, hashSum)


	return hexEncoded
}

// ENDOF function definition
