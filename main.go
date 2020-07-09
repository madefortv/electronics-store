package main

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
