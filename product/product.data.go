package product

import (
	"database/sql"
	"log"

	"github.com/asterisk800/inventoryservice/database"
)

func getProduct(productID int) (*Product, error) {

	row := database.DbConn.QueryRow(`SELECT 
	productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	From products
	WHERE productId = ?`, productID)
	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	// if the raw is empty retrun nil, else return  error
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	log.Printf("Retriving productID: %v", product.ProductID)
	return product, err
}

func removeProduct(productID int) error {
	_, err := database.DbConn.Query(`DELETE FROM products
	WHERE productID = ?`, productID)
	if err != nil {
		return err
	}

	return nil
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT 
		productId,
		manufacturer,
		sku,
		upc,
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
	log.Printf("Retriving %d raws", len(products))
	return products, nil
}

func updateProduct(product Product) error {
	_, err := database.DbConn.Exec(`UPDATE products SET
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=CAST(? AS DECIMAL(13,2)),
	quantityOnHand=?,
	productName=?
	WHERE productID=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID)
	if err != nil {
		return err
	}

	return nil
}

func incertProdcut(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products
	(manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName) VALUES (?,?,?,?,?,?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	log.Printf("Added %v product.", rowsAffected)

	return int(rowsAffected), err
}
