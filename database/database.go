package database

import (
	"database/sql"
	"log"
	"time"
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

	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)

}
