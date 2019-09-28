package models

import (
	"backend/packages/godotenv"
	"backend/packages/gorm"
	_ "backend/packages/gorm/dialects/postgres"
	"fmt"
	"os"
)

var db *gorm.DB //database

func init() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&State{}, &Character{}) //Database migration
	//put in the idnexer to only index upon first launch of binary
	// GetDB().Create(&State{Name: "indexState", Index: true})
}

//GetDB returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
