package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func EstablishConnection() *sql.DB {
	//create a dsn string
	dsn := "user=postgres host=localhost port=5432 dbname=musicdb sslmode=disable"
	// open the connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Cannot create the DB connection: %s", err.Error()))
	}
	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("Cannot ping the DB: %s", err.Error()))
	}

	fmt.Println("Successfully connected to the database!")
	return db
}
