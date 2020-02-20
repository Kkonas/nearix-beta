
package main
 import (
	"log"
	"os"
	"regexp"
	"strings"
	"gopkg.in/yaml.v2"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"encoding/json"
 )
type Pokemon struct{
	Name string
	HashCode string
}
func logErr(err error){
	if(err != nil){
		log.Print(err)
	}
}
func hashes(structC *[]Pokemon){
	output, err := ioutil.ReadFile("hashes.json")
	logErr(err)
	err = json.Unmarshal(output, structC)
	logErr(err)

}
type pokeToAppend struct{
	Name map[string]string
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate){
	if m.Author.ID != "365975655608745985"{
		return
	}else{
		embeds := m.Embeds
		for _, em := range embeds{
			url := em.Image.URL
			var re = regexp.MustCompile(`[0-9]*`)
			replaced := re.ReplaceAllString(em.Title, "")
			finalName := strings.Replace((strings.Replace(strings.Replace(replaced,"Base stats for ","",1),"#","",1)),".","",1)
			hash :=receive(url)
			fmt.Printf(finalName+"\n")
			file, err := os.OpenFile("config/hashes.yaml",os.O_APPEND|os.O_WRONLY, 0600)
			var poke pokeToAppend
			poke.Name = make(map[string]string)
			poke.Name[finalName] = hash
			out , err := yaml.Marshal(&poke)
			defer file.Close()
			_,err = file.Write(out)
			_,err = s.ChannelMessageSend("631206053887082496",hash)
			logErr(err)
		}
	}
}
func main(){
	var pokemon []Pokemon
	hashes(&pokemon)
	client, err := discordgo.New("")
	logErr(err)
	client.AddHandler(messageCreate)
	client.Open()
	for index, name := range pokemon{
		time.Sleep(3 * time.Second)
		fmt.Printf("Now: %s %d/%d\n",name.Name,index,len(pokemon))
		_,err = client.ChannelMessageSend("631206053887082496","p!info "+name.Name)
}
}
