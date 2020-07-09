package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	config         *Config
	productService *ProductService
}

const jsonContentType = "application/json"

func NewServer(config *Config, service *ProductService) *Server {
	return &Server{
		config:         config,
		productService: service,
	}
}

func (s *Server) Handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/products", s.findProducts)
	router.HandleFunc("/products/create", s.createProduct)
	router.HandleFunc("/products/update", s.updateProduct)
	router.HandleFunc("/products/delete", s.deleteProduct)
	router.HandleFunc("/deals", s.deals)
	router.HandleFunc("/offerings", s.offerings)
	router.HandleFunc("/cart", s.cart)
	return router
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}
	httpServer.ListenAndServe()
}

/* Cart Handler */
func (server *Server) cart(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {

	case http.MethodPost:
		var shoppingCart ShoppingCart
		var product Product
		err := json.NewDecoder(request.Body).Decode(&product)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		p := server.productService.getProduct(product)
		if (Product{} == p) {
			http.Error(writer, "Product Does Not Exist", 404)
		}
		err = server.productService.addToCart(p)
		if err != nil {
			http.Error(writer, "Failed to add product to cart", 500)
		}

		items := server.productService.listCartItems()

		if len(items) == 0 {
			shoppingCart = ShoppingCart{}
		} else {
			total, err := server.productService.calculateTotalPrice()
			if err != nil {
				http.Error(writer, "Error calculating total", 500)
			}
			shoppingCart = ShoppingCart{items, total}
		}

		bytes, err := json.Marshal(shoppingCart)
		if err != nil {
			http.Error(writer, "Failed to write a response", 500)
		}

		writer.Header().Set("Content-Type", jsonContentType)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)

	case http.MethodPut:
		var shoppingCart ShoppingCart
		var item Item
		err := json.NewDecoder(request.Body).Decode(&item)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.updateCart(item)
		if err != nil {
			http.Error(writer, "Failed to update cart", 500)
		}

		items := server.productService.listCartItems()

		if len(items) == 0 {
			shoppingCart = ShoppingCart{}
		} else {
			total, err := server.productService.calculateTotalPrice()
			if err != nil {
				http.Error(writer, "Error calculating total", 500)
			}
			shoppingCart = ShoppingCart{items, total}
		}

		bytes, err := json.Marshal(shoppingCart)
		if err != nil {
			http.Error(writer, "Failed to write response", 500)
		}
		writer.Header().Set("Content-Type", jsonContentType)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)

	case http.MethodGet:
		items := server.productService.listCartItems()
		var shoppingCart ShoppingCart

		if len(items) == 0 {
			shoppingCart = ShoppingCart{}
		} else {
			shoppingCart = ShoppingCart{items, "0.0"}
		}

		bytes, err := json.Marshal(shoppingCart)
		if err != nil {
			http.Error(writer, "Failed to write response", 500)
		}
		writer.Header().Set("Content-Type", jsonContentType)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)

	}

}

/* Offerings Handler */
func (server *Server) offerings(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodPost:

		var offering Offering
		err := json.NewDecoder(request.Body).Decode(&offering)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.newOffering(offering)
		if err != nil {
			http.Error(writer, "Failed create new deal", 500)
		}
		writer.WriteHeader(http.StatusCreated)
	}

}

func (server *Server) deals(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:

		deals := server.productService.listDeals()
		bytes, err := json.Marshal(deals)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}
		writer.Header().Set("Content-Type", jsonContentType)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)

	case http.MethodPost:

		var deal Deal
		err := json.NewDecoder(request.Body).Decode(&deal)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.newDeal(deal)
		if err != nil {
			http.Error(writer, "Failed create new deal", 500)
		}
		writer.WriteHeader(http.StatusCreated)
	}

}

func (server *Server) findProducts(writer http.ResponseWriter, request *http.Request) {
	products := server.productService.listProducts()
	bytes, err := json.Marshal(products)
	if err != nil {
		http.Error(writer, "Bad Request", 400)
	}
	writer.Header().Set("Content-Type", jsonContentType)
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func (server *Server) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	var code ProductCode
	err := json.NewDecoder(request.Body).Decode(&code)

	if err != nil {
		http.Error(writer, "Bad Request", 400)
	}

	err = server.productService.deleteProduct(code)
	if err != nil {
		fmt.Errorf("Error in deleting product with id %d, %v", code.Id, err)
		http.Error(writer, "Failed to delete new product", 500)
	}

	writer.WriteHeader(http.StatusOK)
}

func (server *Server) createProduct(writer http.ResponseWriter, request *http.Request) {
	var product Product
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(writer, "Bad Request", 400)
	}
	err = server.productService.newProduct(product)
	if err != nil {
		fmt.Errorf("Error in creating product %q, %v", product, err)
		http.Error(writer, "Failed to create new product", 500)
	}

	writer.WriteHeader(http.StatusCreated)
}

func (server *Server) updateProduct(writer http.ResponseWriter, request *http.Request) {
	var product Product
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		http.Error(writer, "Bad Request", 400)
	}
	err = server.productService.updateProduct(product)
	if err != nil {
		fmt.Errorf("Error in creating product %v", err)
		http.Error(writer, "Failed to update the product", 500)
	}
	writer.WriteHeader(http.StatusNoContent)
}
