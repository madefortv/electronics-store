package store

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestDeals(t *testing.T) {
	// scaffolding
	config := NewConfig()
	productRepository := setupTestDatabase(config)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)

	productService.repository.createDealsTable()

	t.Run("insert new deal", func(t *testing.T) {

		request := newDealRequest(0, 0, "Regular Price", "Retail", "0", "0")
		want := ""

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		var got string
		assertStatus(t, response.Code, http.StatusCreated)
		assertResponseBody(t, got, want)

	})

}

func TestProducts(t *testing.T) {
	// scaffolding
	config := NewConfig()
	productRepository := setupTestDatabase(config)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)

	// database reset seed
	productService.repository.createProductsTable()
	productService.repository.insertProduct(Product{1, "laptop", "very fast", "1000.00"})
	productService.repository.insertProduct(Product{2, "mouse", "much clicky", "10.00"})

	t.Run("get the list of products", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/products", nil)
		want := []Product{{1, "laptop", "very fast", "1000.00"}, {2, "mouse", "much clicky", "10.00"}}

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		var got []Product
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusOK)
		assertProducts(t, got, want)

	})

	t.Run("inserts a new product", func(t *testing.T) {

		request := newCreateProductRequest(0, "monitor", "fourkay", "100.00")
		want := ""
		var got string
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)

		assertResponseBody(t, got, want)

	})

	t.Run("update a product name and description", func(t *testing.T) {

		request := newUpdateProductRequest(1, "laptop", "older", "85.00")
		want := ""

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)

		var got string

		assertStatus(t, response.Code, http.StatusNoContent)
		assertResponseBody(t, got, want)

	})

	t.Run("delete the a product (id = 2)", func(t *testing.T) {
		request := newDeleteProductRequest(2)
		want := ""

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)

		var got string

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)

	})

	t.Run("veirfy the database state", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/products", nil)
		want := []Product{{1, "laptop", "older", "85.00"}, {3, "monitor", "fourkay", "100.00"}}

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		var got []Product
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusOK)
		assertProducts(t, got, want)
	})

}

func setupTestDatabase(config *Config) (repository *ProductRepository) {
	// Helper method for resetting the database
	os.Remove("store.db")
	file, err := os.Create("store.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	db, err := ConnectDatabase(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	repository = NewProductRepository(db)
	return repository
}

func newDealRequest(x, y int, name, coupon, percent string, dtype DealType) *http.Request {
	deal := Deal{
		0,
		name,
		dtype,
		coupon,
		percent,
		x,
		y,
		true,
	}
	body, _ := json.Marshal(deal)
	req, _ := http.NewRequest(http.MethodPost, "/deals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

func newCreateProductRequest(id int, name, description, price string) *http.Request {
	product := Product{
		0,
		name,
		description,
		price,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

type BadRequestBody struct {
	foo string `json:"id"`
}

func buildBadRequest(route string) *http.Request {
	badBody := BadRequestBody{"bar"}
	body, _ := json.Marshal(badBody)
	req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

func newUpdateProductRequest(id int, name, description, price string) *http.Request {
	product := Product{
		id,
		name,
		description,
		price,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPut, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

func newDeleteProductRequest(id int64) *http.Request {
	code := ProductCode{id}
	body, _ := json.Marshal(code)
	req, _ := http.NewRequest(http.MethodDelete, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertProductCode(t *testing.T, got, want ProductCode) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertProduct(t *testing.T, got, want Product) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertProducts(t *testing.T, got, want []Product) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
