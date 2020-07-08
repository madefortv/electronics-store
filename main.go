package main

import (
	"encoding/gob"
	"fmt"
)

func checkNoError(err error, format string) {
	if err != nil {
		panic(fmt.Sprintf(format, err))
	}
}

func init() {
	gob.Register(&Item{})
	gob.Register(&ShoppingCart{})
}

func main() {
	config := NewConfig()

	db, err := ConnectDatabase(config)

	if err != nil {
		panic(err)
	}

	productRepository := NewProductRepository(db)

	productService := NewProductService(config, productRepository)

	server := NewServer(config, productService)

	server.Run()
}
