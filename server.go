package store

import (
	"encoding/json"
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
	router.HandleFunc("/products", s.products)
	router.HandleFunc("/deals", s.deals)
	return router
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}
	httpServer.ListenAndServe()
}

func (server *Server) products(writer http.ResponseWriter, request *http.Request) {
	var product Product
	switch request.Method {
	case http.MethodGet:

		products := server.productService.listProducts()
		bytes, err := json.Marshal(products)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		writer.Header().Set("Content-Type", jsonContentType)
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)

	case http.MethodPost:

		err := json.NewDecoder(request.Body).Decode(&product)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.createProduct(product)
		if err != nil {
			http.Error(writer, "Failed to create new product", 500)
		}

		writer.WriteHeader(http.StatusCreated)

	case http.MethodDelete:
		err := json.NewDecoder(request.Body).Decode(&product)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.repository.deleteProduct(product)
		if err != nil {
			http.Error(writer, "Failed to delete new product", 500)
		}

		writer.WriteHeader(http.StatusOK)

	case http.MethodPut:

		err := json.NewDecoder(request.Body).Decode(&product)
		if err != nil {
			http.Error(writer, "Bad Request", 400)
		}

		err = server.productService.updateProduct(product)
		if err != nil {
			http.Error(writer, "Failed to update the product", 500)
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

/*
func (server *Server) listProducts(writer http.ResponseWriter, request *http.Request) {
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
	var product Product
	err := json.NewDecoder(request.Body).Decode(&product)

	if err != nil {
		http.Error(writer, "Bad Request", 400)
	}

	err = server.productService.repository.deleteProduct(product)
	if err != nil {
		fmt.Errorf("Error in deleting product with id %d, %v", product.Id, err)
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
	err = server.productService.createProduct(product)
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
*/
