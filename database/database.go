package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

func ConectDB() { //ConectDB create a connection with the postgres database
	connStr := "user=postgres dbname=postgres password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err.Error())
	}
	DB = db
}

func CloseDB() { //CloseDB shutdown the connection with the postgres database
	DB.Close()
}
