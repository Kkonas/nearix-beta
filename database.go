
package main
 import (
	"log"
	"regexp"
	"strings"
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
	Name string
	Hash string
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
			file, err := os.OpenFile("goHashes.json:wq
			")
			bytes, err := ioutil.ReadAll()
			var poke pokeToAppend
			err = json.Unmarshal()
			_,err := s.ChannelMessageSend("631206053887082496",hash)
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
	for _, name := range pokemon{
		time.Sleep(3 * time.Second)
		fmt.Printf("Now: "+name.Name+"\n")
		_,err = client.ChannelMessageSend("631206053887082496","p!info "+name.Name)}

}
