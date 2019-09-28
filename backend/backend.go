package main

import (
	"backend/controllers"
	"backend/indexer"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting backend....")
	status := &models.State{}
	models.GetDB().Table("states").Where("name = ?", "indexState").Last(status)
	if status.Index == false {
		fmt.Println("this is the first time this program has been run...")
		fmt.Println("reindexing the API...")
		indexer.Start("https://swapi.co/api/people/")
		fmt.Println("All characters re-index and saved")
		status.Index = true
		models.GetDB().Save(status)
	} else {
		fmt.Println("Already Indexed")
		fmt.Println("The API is Ready")
	}

	router := mux.NewRouter()
	//test the API is working
	router.HandleFunc("/api/check", controllers.Check).Methods("GET")
	//Get the characterdata
	router.HandleFunc("/api/character", controllers.GetCharacter).Methods("POST")
	port := "5555"
	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:5555/api
	if err != nil {
		fmt.Print(err)
	}
}
