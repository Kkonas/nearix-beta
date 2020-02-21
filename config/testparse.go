package main
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)
type PokemonList struct{
	Pokemons map[string]string
}
func main(){
	var str PokemonList
	reader, err := ioutil.ReadFile("hashes.yaml")
	if err != nil{
		log.Fatal(err)
	}
	yaml.Unmarshal(reader,str)
	fmt.Printf(str.Pokemons["Bulbasaur"])
}
