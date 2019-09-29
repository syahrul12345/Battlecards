//Package models represents all character models that is received from SWAPI, the postgres character model, and the http response model.
package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//State is the flag to only index when first starting
type State struct {
	gorm.Model
	Name  string
	Index bool
}

//CharacterResponse is the JSON object of a character when it is first returend by
type CharacterResponse struct {
	Name      string
	Gender    string
	StarShips []string
	Vehicles  []string
	HomeWorld string
}

//Character represents the JSON object of a StarWars character after reindexing
type Character struct {
	gorm.Model
	NameSearch string
	Name       string
	Gender     string
	StarShips  pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Vehicles   pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Home       pq.ByteaArray `gorm:"type:varchar(8000)[]"`
}

//CharacterResult is the JSON object send to the front end
type CharacterResult struct {
	Name      string
	Gender    string
	StarShips []StarShip
	Vehicles  []Vehicle
	HomeWorld Home
}

//Cache represents the data that will be written to the local cache as a text file
type Cache struct {
	Time       int64
	Characters []CharacterResult
}

//StarShip represents the JSON object of a one StarShip after reindexing
type StarShip struct {
	Model             string
	Starship_Class    string
	Hyperdrive_Rating string
	Cost_In_Credits   string
	Manufacturer      string
}

//Vehicle represents the JSON object of one Vehicle after reindexing
type Vehicle struct {
	Name          string
	Model         string
	CostInCredits string
}

//Home represents the homeworld stats
type Home struct {
	Name       string
	Population string
	Climate    string
}

//Cache will save the current character result to a local text file
//This function is implemented on a *CharacterResult
func (charRes *CharacterResult) Cache() {
	fmt.Println("########################################################")
	//1) Open the cache file in ../cache/cache.txt
	//2) read the data inside
	//2) Parse the data into a Cache struct object
	//3) Append the current characterResult to the Cache Object, if it doesn't already exist
	//4) Save the cache object back to the cache file
	//5) Updates the time of the cache
	//Open the File
	file, fileErr := os.OpenFile("../cache/Cache.txt", os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if fileErr != nil {
		//File doesnt exist, so we need to create it
		fmt.Println("Cache file doesnt exist")
		fmt.Println("Creating cache.txt in /cache/ folder...")
		file, _ = os.Create("../cache/Cache.txt")
		fmt.Println("File Created")
	}
	defer file.Close()
	//Read the file
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to read the data from cache file, possibly not bytes")
	}
	cache := &Cache{}
	bodyErr := json.Unmarshal(b, cache)
	if bodyErr != nil {
		if bodyErr.Error() == "unexpected end of JSON input" {
			fmt.Println("Cache file is empty...")
		}
	}
	//set the duplicate flag.
	duplicate := false
	for _, character := range cache.Characters {
		if character.Name == charRes.Name {
			fmt.Println("This search already exists in the cache")
			duplicate = true
		}
	}
	//Only if it isn't a dupicate do we add to the cache file
	if duplicate != true {
		//Add the current characters to be saved
		cache.Time = time.Now().Unix()
		cache.Characters = append(cache.Characters, *charRes)
		//Let's marshal it back into bytes, for saving
		cacheBytes, cacheErr := json.Marshal(cache)
		if cacheErr != nil {
			fmt.Println("Added response to Go Cache object, but could not conver it to bytes")
		}

		//Overwrites the older Cache struct with the newer Cache struct.
		//The index 0 in the second argument will write the new Cache struct from index 0 in the file thus overwrite
		_, writeErr := file.WriteAt(cacheBytes, 0)
		if writeErr != nil {
			fmt.Println(writeErr)
			fmt.Println("Failed to save cache data to file")
		}
		fmt.Printf("Character %s has been saved to cache\n", string(charRes.Name))
	}
	fmt.Println("########################################################")
}

//ToLowerAndNoSpecial removes all whitespaces, diacratics and special charates such as / and -
func ToLowerAndNoSpecial(str string) string {
	return strings.Map(func(r rune) rune {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return -1
		}
		return r
	}, normalize(strings.ToLower(str)))
}

//Helper function for ToLowerAndNoSpecial
func normalize(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}

//Helper function for ToLowerAndNoSpecial
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
