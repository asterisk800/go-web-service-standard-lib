package database

import (
	"database/sql"
	"log"
)

// DbConn : capital because we want to export it and usee it in other modules.
var DbConn *sql.DB

// SetupDatabase :
func SetupDatabase() {
	var err error

	DbConn, err = sql.Open("mysql", "inventory:inventory@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Connected to DB within %v ns", DbConn.Stats().WaitDuration.Nanoseconds())
}
