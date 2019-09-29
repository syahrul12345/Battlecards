package main

import (
	"backend/controllers"
	"backend/indexer"
	"backend/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	//API
	//test the API is working
	router.HandleFunc("/api/check", controllers.Check).Methods("GET")
	//Get the characterdata, and immediately stores it to cache
	router.HandleFunc("/api/character", controllers.GetCharacter).Methods("POST")
	//Get the cacheData
	router.HandleFunc("/api/getCache", controllers.GetCharacterCache).Methods("GET")

	//Serve the Compiled VUE.js frontend
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../dist/")))
	port := "5555"
	fmt.Println("Serving static website at http://localhost:5555/")
	//lets set the cors policy for testing
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	err := http.ListenAndServe(":"+port, handler) //Launch the app, visit localhost:5555/api
	if err != nil {
		fmt.Print(err)
	}
}
