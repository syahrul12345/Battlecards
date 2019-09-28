package controllers

import (
	"backend/models"

	"backend/utils"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

//RequestPayload is the incoming request object
type RequestPayload struct {
	Name string
}

//GetCharacter gets the corresponding character JSON from the database.
var GetCharacter = func(writer http.ResponseWriter, request *http.Request) {
	requestPayload := &RequestPayload{}
	err := json.NewDecoder(request.Body).Decode(requestPayload)
	if err != nil {
		utils.Respond(writer, utils.Message(false, "Incorrect request format"))
	}
	if requestPayload.Name == "" {
		utils.Respond(writer, utils.Message(false, "Missing data in the name parameter of request"))
	} else {
		name := requestPayload.Name
		character := &models.Character{}
		err := models.GetDB().Table("characters").Where("name_search = ?", models.ToLowerAndNoSpecial(name)).First(character).Error
		//Handle errors if we cannot get the database
		if err != nil {
			//Record doesnt exist
			if err == gorm.ErrRecordNotFound {
				utils.Respond(writer, utils.Message(false, "No such character with that name"))
			} else {
				//Database network errors
				utils.Respond(writer, utils.Message(false, "Database Error"))
			}
		} else {
			//character.Vehicles and character.Starship exists in a pq.ByteArray format when it is returned from postgreSQL
			//We would like to make them into the correct json format
			vehicles := getVehicleFromBytes(character.Vehicles)
			starships := getStarshipFromBytes(character.StarShips)
			home := getHomeFromBytes(character.Home)
			//build the response to send to the frontend
			resp := &models.CharacterResult{Name: character.Name, Gender: character.Gender, HomeWorld: home}
			if len(vehicles) != 0 {
				resp.Vehicles = vehicles
			}
			if len(starships) != 0 {
				resp.StarShips = starships
			}
			utils.Respond(writer, utils.Message(true, resp))
		}

	}
}

//Check the API is working
var Check = func(writer http.ResponseWriter, request *http.Request) {
	utils.Respond(writer, utils.Message(true, "API is WORKING"))
}

//Accepts the bytearray of all Starships returned from postgress
func getStarshipFromBytes(starshipArray pq.ByteaArray) []models.StarShip {
	starships := []models.StarShip{}
	for _, item := range starshipArray {
		starship := &models.StarShip{}
		json.Unmarshal(item, starship)
		starships = append(starships, *starship)
	}
	return starships
}

//Accepts the bytearray of all Vehicles returned from postgress
func getVehicleFromBytes(vehicleArray pq.ByteaArray) []models.Vehicle {
	vehicles := []models.Vehicle{}
	for _, item := range vehicleArray {
		vehicle := &models.Vehicle{}
		json.Unmarshal(item, vehicle)
		vehicles = append(vehicles, *vehicle)
	}
	return vehicles
}

//Accets the bytearray of homeworld from postgres
func getHomeFromBytes(homeArray pq.ByteaArray) models.Home {
	home := &models.Home{}
	for _, item := range homeArray {
		err := json.Unmarshal(item, home)
		if err != nil {
			utils.Message(false, err.Error())
		}
	}
	return *home
}
