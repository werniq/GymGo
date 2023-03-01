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

var (
	tables = [7]string{
		"legs", "chest", "glutes", "back", "biceps", "triceps", "shoulders",
	}
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

// Serve function configures/connects all routes and handlers
func (web *webapp) Serve() error {
	router := gin.Default()

	router.Use(web.CorsMiddleware())

	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error loading dotenv file: %v", err)
	}

	api := router.Group("/api")

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
	// | Regenerate				|							|
	// |						|							|
	// | Dumbbell Incline Press	|		Video 				|
	// | Technique				|							|
	// | Regenerate				|							|
	// |						|							|
	// |----------------------------------------------------|

	rearUpper := api.Group("/rear-upper")
	{
		rearUpper.POST("/biceps", web.ReturnAllBicepsExercises)
		rearUpper.POST("/shoulders", web.ReturnAllShouldersExercises)
		rearUpper.POST("/back", web.ReturnAllBackExercises)
	}

	frontUpper := api.Group("/front-upper")
	{
		frontUpper.POST("/triceps", web.ReturnAllTricepsExercises)
		frontUpper.POST("/chest", web.ReturnAllChestExercises)
		frontUpper.POST("/shoulders", web.ReturnAllShouldersExercises)
	}

	legs := api.Group("/legs")
	{
		legs.POST("/legs", web.ReturnAllLegsExercises)
		legs.POST("/glutes", web.ReturnAllGlutesExercises)
	}

	api.POST("/generate-workout", web.GenerateWorkout)

	return router.Run(":4001")
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}
	dsn := "user=postgres password=Matwyenko1_ dbname=workout-website host=localhost port=5432 sslmode=disable"
	portSTR := "localhost:4001//"

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

	// legs, chest, back, biceps, glutes exercises already in databases
	if err := web.Serve(); err != nil {
		fmt.Printf("Error serving connection: %v", err)
	}
}
