package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/werniq/gym/driver"
	"github.com/werniq/gym/models"
	"html/template"
	"log"
	"os"
	"strconv"
)

type webapp struct {
	env    string
	url    string
	port   int
	stripe struct {
		secret string
		key    string
	}
	db            models.DatabaseModel
	dbDSN         string
	errorLog      *log.Logger
	templateCache map[string]*template.Template
}

func (web *webapp) Serve() error {
	router := gin.Default()

	//api := router.Group("/muscles")

	router.GET("/home", web.HomePage)

	// need to create following tables, and insert values in such way:
	// video for quick watch
	// name of exercise
	// description/technique

	// left side of website -> workout
	// right side -> info about exercises (as shown above)
	// workout will be generated in such way:
	// user chooses, which muscles wants to train
	// (optional, chooses difficulty)
	// chooses for each of exercises in given area
	// for example
	// chest 2
	// barbell press
	// dumbbell incline press
	// and video for quick watch, and technique description
	// |----------------------------------------------------|
	// | NavBar section										|
	// |====================================================|
	// | Barbell Press			|		Video				|
	// | Technique  			|							|
	// | Dumbbell Incline Press	|		Video 				|
	// | Technique				|							|
	// |						|							|
	// |						|							|
	// |----------------------------------------------------|
	//rearUpper := api.Group("/rear-upper")
	//{
	//	rearUpper.GET("/biceps")
	//	rearUpper.GET("/shoulders")
	//	rearUpper.GET("/back")
	//}
	//
	//frontUpper := api.Group("/front-upper")
	//{
	//	frontUpper.GET("/triceps")
	//	frontUpper.GET("/chest")
	//	frontUpper.GET("/shoulders")
	//}
	//
	//legs := api.Group("/legs")
	//{
	//	legs.GET("/glutes")
	//	legs.GET("/quads")
	//	legs.GET("/hamstrings")
	//}

	return router.Run(":8000")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	dsn := os.Getenv("DB_DSN")
	portSTR := os.Getenv("HOST")

	db, err := driver.OpenDB()
	if err != nil {
		fmt.Printf("Error retrieving database connection: %v", err)
	}

	portINT, _ := strconv.Atoi(portSTR)
	templateCache := make(map[string]*template.Template)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Lshortfile|log.Ldate|log.Ltime)

	web := webapp{
		env:  "development",
		url:  fmt.Sprintf("http://localhost:%d", portINT),
		port: portINT,
		stripe: struct {
			secret string
			key    string
		}{},
		db: models.DatabaseModel{
			DB: db,
		},
		errorLog:      errorLog,
		dbDSN:         dsn,
		templateCache: templateCache,
	}

	router := gin.Default()

	router.GET("/", web.HomePage)
	router.GET("/legs", web.Legs)
	router.GET("/chets", web.Chest)
	router.GET("/back", web.Back)
	router.GET("/biceps", web.Biceps)

	router.Run(":8000")
}
