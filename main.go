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
		fmt.Printf(Lang(confStruct,langStruct,1,"emptytoken"))
		response := readStdin()
		confStruct.Token = response
	}
}
func Lang(confStruct *Config,langStruct *LangConfig, id int,message string) string{
	return langStruct.Languages[confStruct.Constants.Language][id][message]
}
func logErr(err error){
	if err != nil{
		log.Print(err)
	}
}
func readStdin() string{
	reader := bufio.NewReader(os.Stdin)
	raw ,_ := reader.ReadString('\n')
	return strings.Replace(raw,"\n","",1)
}
/*func messageCreate(session *discordgo.Session,message *discordgo.MessageCreate){
	if message.Author.ID != conf.Constants.PokeCordID{
	return
	}else{
		embeds := message.Embeds
		for _, embed := range embeds{
			spawnUrl := embed.Image.URL

		}
	}
}*/
// ENDOF function declaration section
// Begin structs
type LangConfig struct{
	Languages	map[string][]map[string]string
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
	var lang LangConfig
	var conf Config
	initCheck(&conf,&lang)
	fmt.Printf("Token: %s, Version: %s, Guilds: %s ,ID: %s\n",conf.Token,conf.Version, conf.Session.Guilds, conf.Constants.PokeCordID)
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
	writeConfigYaml(configFile,&conf)
}
