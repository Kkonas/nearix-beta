/* Welcome to the discord multi spam bot
 Author: SteeW (a.k.a joker-ware)
 Following code may or may not be commented.
 "Abandon all hope ye who enter here" */

package main
import (
	"fmt"
	"bufio"
	"strings"
	"io/ioutil"
	"os"
	"log"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)
// Begin CONSTS

const configFile = "config/config.yaml"

// ENDOF CONSTS
// Begin function declaration section
func readFile(path string) []byte{
	file, err := ioutil.ReadFile(path)
	logErr(err)
	return file
}
func writeFile(path string, buffer []byte){
	err := ioutil.WriteFile(path, buffer, 0200)
	logErr(err)
}
func readYaml(path string, confStruct *Config){
	reader := readFile(path)
	err := yaml.Unmarshal(reader, confStruct)
	logErr(err)
}
func writeYaml(path string,confStruct *Config){
	output, err := yaml.Marshal(confStruct)
	logErr(err)
	writeFile(path, output)
}
func initCheck(confStruct *Config){
	fmt.Printf("Do you wanna restore last session? Y/n:\n")
	response := strings.ToLower(readStdin())
	if response == "y"{
		readYaml(configFile, confStruct)
	}else if response == "n"{
		return
	}else{
		os.Exit(3)
	}
}
func logErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}
func readStdin() string{
	reader := bufio.NewReader(os.Stdin)
	raw ,_ := reader.ReadString('\n')
	return strings.Replace(raw,"\n","",1)
}
// ENDOF function declaration section
// Begin structs
type Config struct{
	Token		string
	Version		string
	Session		State
}
type State struct{
	Guilds		[]Guild
}
type Guild struct{
	Id		string
	Name		string
	Channels	[]Channel
}
type Channel struct{
	Id		string
	Name		string
	Enabled		bool
}
// ENDOF structs

// Begin main
func main(){
	var conf Config
	initCheck(&conf)
	fmt.Printf("Token: %s, Version: %s, Guilds: %s \n",conf.Token,conf.Version, conf.Session.Guilds)
	client, err := discordgo.New(conf.Token)
	logErr(err)
	err = client.Open()
	logErr(err)
	for index, guild := range client.State.Guilds{
		appendGuild := Guild{guild.ID,guild.Name,nil}
		conf.Session.Guilds = append(conf.Session.Guilds,appendGuild)
		for _, channel := range guild.Channels{
			appendChannel := Channel{channel.ID,channel.Name,false}
			conf.Session.Guilds[index].Channels = append(conf.Session.Guilds[index].Channels, appendChannel)
		}
	}
	writeYaml(configFile,&conf)
}
