// Welcome to the discord multi spam bot
// Author: SteeW (a.k.a joker-ware)
// Following code may or may not be commented.
// "Abandon all hope ye who enter here"

package main
import (
	"fmt"
	"bufio"
	"strings"
	"io/ioutil"
	"os"
	"log"
	//"path/filepath"
	"github.com/bwmarrin/discordgo"
	//"gopkg.in/yaml.v2"
)
// Begin function declaration section
func readFile(path string){
	raw, _ := ioutil.ReadFile(string(path))
	fmt.Printf(string(raw))
}
func init(){
	// Perform initial check

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
	token string
	version string
}
type State struct{
	stateToken string
	stateServers struct{
		id string
		stateChannels []string
	}
}
// ENDOF structs
// Begin main
func main(){
	client, err := discordgo.New("NjIyNTEzNTIwMzc3ODU2MDEw.XkMDww.US7_wwkIbf6luj3Pcd0bCj8I0HA")
	logErr(err)
	err = client.Open()
	logErr(err)
	for _, guild := range client.State.Guilds{
		fmt.Printf(guild.ID+"\n")
	}
}
