package driver

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/locales/dyo_SN"
	_ "github.com/lib/pq"
)

const (
	legExercises = 15
)

func OpenDB() (*sql.DB, error) {
	dsn := "user=postgres password=Matwyenko1_ dbname=workout-website host=localhost port=5432 sslmode=disable"

	fmt.Printf("Database dsn: %s", dyo_SN.New())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Error pinging database connection: %v", err)
		return nil, err
	}

	return db, nil
}
