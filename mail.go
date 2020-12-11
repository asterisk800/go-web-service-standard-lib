package main

import (
	"net/http"

	"github.com/asterisk800/inventoryservice/product"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe("localhost:5000", nil)

}
