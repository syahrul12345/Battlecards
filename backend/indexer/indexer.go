//Package indexer will reindex the SWAPI, so it becomes searchable by name. It is then stored in a local postgresSQL database
package indexer

import (
	"backend/models"
	"backend/requests"
	"backend/utils"
	"encoding/json"
	"fmt"
)

//ParsedResponse returned by the SWAPI when calling /api/people.
type ParsedResponse struct {
	Next    string
	Results []models.CharacterResponse
}

//Start will start the indexer feature of the backend. Indexer makes an API call to the SWAPI, and then re-indexes the data
//It is impossible to search for character names in SWAPI, hence a the dataset has to be reindexed by character name.
func Start(rpc string) {
	next := save(rpc)
	for next != "" {
		next = save(next)
	}
}
func save(rpc string) string {
	//Gets the response from the SWAPI
	response := requests.Get(rpc)
	//parse the response into a response object
	var parsedResponse = new(ParsedResponse)
	parsedError := json.Unmarshal(toByte(response["message"]), &parsedResponse)
	if parsedError != nil {
		utils.Message(false, parsedError.Error())
	}
	for i := range parsedResponse.Results {
		currentCharacter := &parsedResponse.Results[i]
		go saveCharacter(currentCharacter)

	}
	return parsedResponse.Next
}

func saveCharacter(characterResponse *models.CharacterResponse) {
	//Create the  Character object which will be saved to postgresSQL.
	//We assign the values which Character and CharacterResponse both share, which is Name and Gender

	//Since we want to search by name, let's store the Name in lowercase, and remove all special features
	// so if a user searches for R2-D2 -> the backend should search for r2d2
	character := &models.Character{NameSearch: models.ToLowerAndNoSpecial(characterResponse.Name), Name: characterResponse.Name, Gender: characterResponse.Gender}
	//We want to parse CharacterResponse into Character
	//CharacterResponse has field StarShips and Vehicles but these are links, we need to
	//replace each link with a Starship and Vehicle json
	//need to handle if the length of the StarShips array of links is 0

	//Channels for Asynchronous code
	starShipChan := make(chan *models.StarShip)
	vehicleChan := make(chan *models.Vehicle)
	homeChan := make(chan *models.Home)
	go func() {
		if len(characterResponse.StarShips) != 0 {
			//For performane, we will make RPC calls to starship & vehicle links asyncronously. character.StarShips[i]
			// represents a starship link for example https://swapi.co/api/starship/30/
			// The link is then used as input to getStarships(link), which will return the JSON object of the starship
			// Eg: getStarships('https://swapi.co/api/starship/30/) gives Starship = {Model:"T-65 X-Wing",Class:"Starfighter",hyperdrive_rating:"1.0"}
			// We will then add this to the Character, which will be saved to the local DB.
			for i := range characterResponse.StarShips {
				//execute asyncrhonously and store it in a channel
				go func(i int, co chan<- *models.StarShip) {
					co <- getStarships(characterResponse.StarShips[i])
				}(i, starShipChan)
			}
		}
	}()
	go func() {
		//Do exactly the same for the Vehicle array of links
		if len(characterResponse.Vehicles) != 0 {
			for i := range characterResponse.Vehicles {
				go func(i int, co chan<- *models.Vehicle) {
					co <- getVehicles(characterResponse.Vehicles[i])
				}(i, vehicleChan)
			}

		}
	}()
	//We will then parse homeworld information for the Character. CharacterResponse.HomeWorld returns the string of the homworld if any
	go func(co chan<- *models.Home) {
		if len(characterResponse.HomeWorld) != 0 {
			co <- getHome(characterResponse.HomeWorld)
		}
	}(homeChan)
	//Finally we append it to the StarShip list into the Character json that will be saved to the server
	for range characterResponse.StarShips {
		//Since postgresSQL is incapable of storing array of structs, we convert all nested structs to bytearrays
		starShipRaw, _ := json.Marshal(<-starShipChan)
		character.StarShips = append(character.StarShips, starShipRaw)
	}
	for range characterResponse.Vehicles {
		vehicleRaw, _ := json.Marshal(<-vehicleChan)
		character.Vehicles = append(character.Vehicles, vehicleRaw)
	}
	//We'll need to marshal the Home struct into a bytearray for postgreSQL
	homeRaw, _ := json.Marshal(<-homeChan)
	character.Home = append(character.Home, homeRaw)
	models.GetDB().Create(character)
	fmt.Println("saved for: ", character.Name)

}

//get the starship based on the link, and then returns the StarShip JSON
func getStarships(rpc string) *models.StarShip {
	response := requests.Get(rpc)
	starShip := &models.StarShip{}
	parsedError := json.Unmarshal(toByte(response["message"]), starShip)
	if parsedError != nil {
		fmt.Println("Failed to get starship")
	}
	return starShip
}

//get the vehicle based on the link, and then returns the Vehicle JSON
func getVehicles(rpc string) *models.Vehicle {
	response := requests.Get(rpc)
	vehicle := &models.Vehicle{}
	parsedError := json.Unmarshal(toByte(response["message"]), vehicle)
	if parsedError != nil {
		fmt.Println("failed to get starship")
	}
	return vehicle
}
func getHome(rpc string) *models.Home {
	response := requests.Get(rpc)
	home := &models.Home{}
	parsedError := json.Unmarshal(toByte(response["message"]), home)
	if parsedError != nil {
		fmt.Println("Faeild to get home")
	}
	return home

}

func toByte(inter interface{}) []byte {
	return inter.([]byte)
}
