package main

import (
	"log"
	"net/http"

	"github.com/asterisk800/inventoryservice/database"
	"github.com/asterisk800/inventoryservice/product"

	// we are using _ just to import the drive for it side effect to not to explicitly improt the driver.
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	err := http.ListenAndServe("localhost:5000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
