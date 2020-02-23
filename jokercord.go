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
	"os/signal"
	"syscall"
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
func readConfigYaml(path string, structPointer *Config){
	reader := readFile(path)
	err := yaml.Unmarshal(reader, structPointer)
	logErr(err)
}
func readLangYaml(path string, structPointer *LangConfig){
	reader := readFile(path)
	err := yaml.Unmarshal(reader, structPointer)
	logErr(err)
}

func writeConfigYaml(path string,confStruct *Config){
	output, err := yaml.Marshal(confStruct)
	logErr(err)
writeFile(path, output)
}
func writeLangYaml(path string,langStruct *LangConfig){
	output, err := yaml.Marshal(langStruct)
	logErr(err)
writeFile(path, output)
}

func initCheck(confStruct *Config,langStruct *LangConfig){
	readConfigYaml("config/config.yaml",confStruct)
	readLangYaml("config/languages.yaml",langStruct)
	if confStruct.Token == ""{
		fmt.Printf(Lang("emptytoken"))
		fmt.Printf(Lang("tokenprompt"))
		response := readStdin()
		confStruct.Token = response
		updateConfigYaml()
	}
}
func Lang(message string) string{
	return lang.Languages[conf.Constants.Language][message]
}
func logErr(err error){
	if err != nil{
		readConfigYaml("config/config.yaml",&conf)
		readLangYaml("config/languages.yaml",&lang)
		log.Print(lang.Languages[conf.Constants.Language]["error"])
	}
}
func readStdin() string{
	reader := bufio.NewReader(os.Stdin)
	raw ,_ := reader.ReadString('\n')
	return strings.Replace(raw,"\n","",1)
}
func messageCreate(session *discordgo.Session,message *discordgo.MessageCreate){
	if message.Author.ID != conf.Constants.PokeCordID{
	return
	}else{
		if message.Embeds == nil{
			return
		}else{
		embeds := message.Embeds
		for _, embed := range embeds{
			if (embed.Image == nil){
				return
			}else{
				spawnUrl := embed.Image.URL
				session.ChannelMessageSend(message.ChannelID,"p!catch "+receive(spawnUrl))
				}
			}
		}
	}
}
func refresh(client *discordgo.Session){
	if len(conf.Session.Guilds) == 0{
	for index, guild := range client.State.Guilds{
		appendGuild := Guild{guild.ID,guild.Name,nil}
		conf.Session.Guilds = append(conf.Session.Guilds,appendGuild)
		for _, channel := range guild.Channels{
			appendChannel := Channel{channel.ID,channel.Name,false}
			conf.Session.Guilds[index].Channels = append(conf.Session.Guilds[index].Channels, appendChannel)
		}
	}
	if len(conf.Session.Guilds) == 0{
		return
	}else{
	updateConfigYaml()
	}
	}
}
func updateConfigYaml(){
	writeConfigYaml("config/config.yaml",&conf)
}
func init(){
	initCheck(&conf,&lang)
	readConfigYaml(configFile,&conf)
	client, err := discordgo.New(conf.Token)
	logErr(err)
	client.AddHandler(messageCreate)
	if conf.Constants.First == true{
		fmt.Printf("Please choose your language code. Available languages are: ")
		for languageName,_ := range lang.Languages{
			fmt.Printf(languageName+"|")
		}
		fmt.Printf(" :")
		selectedLang := readStdin()
		if lang.Languages[selectedLang] != nil || len(lang.Languages[selectedLang]) != 0{
			conf.Constants.Language = selectedLang
			conf.Constants.First = false
			updateConfigYaml()
		}else{
			fmt.Printf("Chosen language is not valid.")
			log.Fatal("Language error.")
		}
		err = client.Open()
		if err != nil{
		log.Fatal(lang.Languages[conf.Constants.Language]["tokenerror"])}
		refresh(client)
	}else{
		err = client.Open()
	}
}
// ENDOF function declaration section
// Begin structs
type LangConfig struct{
	Languages	map[string]map[string]string
}
type Config struct{
	Token		string
	Constants	Const
	Version		string
	Session		State
}
type Const struct{
	PokeCordID	string
	Language	string
	First		bool
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
// GLOBAL values
var lang LangConfig
var conf Config
var client discordgo.Session
// ENDOF GLOBAL values
func main(){
	fmt.Printf("Token: %s, Version: %s, ID: %s\n",conf.Token,conf.Version, conf.Constants.PokeCordID)

	fmt.Println(lang.Languages[conf.Constants.Language]["running"])
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	client.Close()
}
