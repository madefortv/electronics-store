package store

import (
	"fmt"
	"log"
	"os"
)

func checkNoError(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err))
	}
}

func main() {
	config := NewConfig()
	db, err := ConnectDatabase(config)

	if err != nil {
		panic(err)
	}

	productRepository := NewProductRepository(db)

	productService := NewProductService(config, productRepository)

	os.Remove("store.db")
	file, err := os.Create("store.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	db, err = ConnectDatabase(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	productService.repository.createProductsTable()
	productService.repository.createDealsTable()
	productService.repository.createOfferingTable()

	server := NewServer(config, productService)

	server.Run()
}
