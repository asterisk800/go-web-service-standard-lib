package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/asterisk800/inventoryservice/database"
)

// Most of the code in this file is used to emulate how we will be working with a database

// used to hold our product list in memory in a map
var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}
	return prodMap, nil
}

func getProduct(productID int) (*Product, error) {

	row := database.DbConn.QueryRow(`SELECT 
	productId,
	manufacturer,
	sku.
	upc.
	pricePerUnit,
	quantityOnHand,
	productName
	From products,
	WHERE productId=?`, productID)
	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand)
	// if the raw is empty retrun nil, else return  error
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return product, err
}

func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT 
		productId,
		manufacturer,
		sku.
		upc.
		pricePerUnit,
		quantityOnHand,
		productName
		From products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}
	log.Printf("Retriving result: %v", products)
	return products, nil
}

func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductID() int {
	productIds := getProductIds()
	return productIds[len(productIds)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	// if the product id is set, update, otherwise add
	addOrUpdateID := -1
	if product.ProductID > 0 {
		oldProduct, err := getProduct(product.ProductID)
		if err != nil {
			return addOrUpdateID, err
		}
		// if it exists, replace it, otherwise return error
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
