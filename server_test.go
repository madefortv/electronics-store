package main

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

func TestShoppingCart(t *testing.T) {
	// scaffolding
	config := NewConfig()
	productRepository := setupTestDatabase(config)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)
	//add cart table
	productService.repository.createCartTable()

	// some deals to offer
	productService.repository.createDealsTable()
	productService.repository.insertDeal(Deal{Name: "Regular Price", Type: "Retail"})
	productService.repository.insertDeal(Deal{Name: "Half Off", Type: "Percent", Percent: "50"})
	productService.repository.insertDeal(Deal{Name: "Laptop Mouse Bundle", Type: "Bundle"})
	productService.repository.insertDeal(Deal{Name: "Buy 3 USB get 1 free", Type: "BuyXGetYFree", X: 3, Y: 1})

	// some products to list
	productService.repository.createProductsTable()
	productService.repository.insertProduct(Product{1, "laptop", "very fast", "1000.00"})
	productService.repository.insertProduct(Product{2, "mouse", "much clicky", "10.00"})
	productService.repository.insertProduct(Product{3, "monitor", "four kay", "100.00"})
	productService.repository.insertProduct(Product{4, "usb", "type see", "1.00"})

	// actual items
	productService.repository.createOfferingsTable()
	// regular priced mouse
	productService.repository.insertOffering(Offering{ProductId: 2, DealId: 1})
	// laptop with a mouse free
	productService.repository.insertOffering(Offering{ProductId: 2, DealId: 3})
	productService.repository.insertOffering(Offering{ProductId: 1, DealId: 3})
	// 50% off monitors
	productService.repository.insertOffering(Offering{ProductId: 3, DealId: 2})

	t.Run("get empty cart", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
		req.Header.Set("Content-Type", jsonContentType)

		want := ShoppingCart{}

		var got ShoppingCart
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

		assertShoppingCart(t, got, want)
	})

	t.Run("Add an item to a shopping cart", func(t *testing.T) {

		body, _ := json.Marshal(Product{Id: 3, Name: "monitor", Description: "four kay", Price: "100.00"})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{Id: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 1}}
		want := ShoppingCart{Items: items, Total: "100.00"}

		var got ShoppingCart
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

		assertShoppingCart(t, got, want)

	})

	t.Run("Modify the quantity of a certain product", func(t *testing.T) {

		body, _ := json.Marshal(Item{Product{3, "monitor", "four kay", "100.00"}, 2})

		req, _ := http.NewRequest(http.MethodPut, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{3, "monitor", "four kay", "100.00"}, 2}}
		want := ShoppingCart{Items: items, Total: "200.00"}

		var got ShoppingCart
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)
		assertShoppingCart(t, got, want)
	})

}

func TestOfferings(t *testing.T) {
	// scaffolding
	config := NewConfig()
	productRepository := setupTestDatabase(config)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)

	// some deals to offer
	productService.repository.createDealsTable()
	productService.repository.insertDeal(Deal{Name: "Regular Price", Type: "Retail"})
	productService.repository.insertDeal(Deal{Name: "Half Off", Type: "Percent", Percent: "50"})
	productService.repository.insertDeal(Deal{Name: "Laptop Mouse Bundle", Type: "Bundle"})
	productService.repository.insertDeal(Deal{Name: "Buy 3 USB get 1 free", Type: "BuyXGetYFree", X: 3, Y: 1})

	// some products to list
	productService.repository.createProductsTable()
	productService.repository.insertProduct(Product{1, "laptop", "very fast", "1000.00"})
	productService.repository.insertProduct(Product{2, "mouse", "much clicky", "10.00"})
	productService.repository.insertProduct(Product{3, "monitor", "four kay", "100.00"})
	productService.repository.insertProduct(Product{4, "usb", "type see", "1.00"})

	// actual items
	productService.repository.createOfferingsTable()
	// regular priced mouse
	productService.repository.insertOffering(Offering{ProductId: 2, DealId: 1})
	// laptop with a mouse free
	productService.repository.insertOffering(Offering{ProductId: 2, DealId: 3})
	productService.repository.insertOffering(Offering{ProductId: 1, DealId: 3})
	// 50% off monitors
	productService.repository.insertOffering(Offering{ProductId: 3, DealId: 2})

	t.Run("create new offering connecting usbs to the buy 3 USBs get 1 free offering", func(t *testing.T) {

		body, _ := json.Marshal(Offering{ProductId: 4, DealId: 4})
		req, _ := http.NewRequest(http.MethodPost, "/offerings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		want := ""
		var got string
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

		assertStatus(t, response.Code, http.StatusCreated)

		assertResponseBody(t, got, want)
	})

	// Get all active offers associated with a specific product and all other products associated with that offer/deal
	/*
	   SELECT products.name, products.price, deals.name, deals.type FROM offerings WHERE active = 1 INNER JOIN products on products.id = offerings.product_id INNER JOIN deals on deals.id = offerings.deal_id;

	   TODO: restrict to a specific product/deal?
	*/

}

func TestDeals(t *testing.T) {
	// scaffolding
	config := NewConfig()
	productRepository := setupTestDatabase(config)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)

	// database reset seed
	productService.repository.createDealsTable()

	productService.repository.insertDeal(Deal{Name: "Regular Price", Type: "Retail"})
	productService.repository.insertDeal(Deal{Name: "Half Off", Type: "Percent", Percent: "50"})
	//productService.repository.insertDeal(Deal{Name: "Buy One Get One Free", Type: "BuyXGetYFree", X: 1, Y: 1})
	//productService.repository.insertDeal(Deal{Name: "Get a Mouse free with any Laptop", Type: "Bundle"})

	t.Run("get the list of deals", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/deals", nil)
		want := []Deal{{Id: 1, Name: "Regular Price", Type: "Retail"}, {Id: 2, Name: "Half Off", Type: "Percent", Percent: "50"}}

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		var got []Deal
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", response.Body, err)
		}
		assertStatus(t, response.Code, http.StatusOK)
		assertDeals(t, got, want)

	})
	t.Run("inserts a new deal", func(t *testing.T) {

		body, _ := json.Marshal(Deal{Name: "Half off any regular price item", Type: "Percent", Percent: "50"})
		req, _ := http.NewRequest(http.MethodPost, "/deals", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		want := ""
		var got string
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

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

func newCreateProductRequest(id int, name, description, price string) *http.Request {
	product := Product{
		0,
		name,
		description,
		price,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPost, "/products/create", bytes.NewBuffer(body))
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
	req, _ := http.NewRequest(http.MethodPost, "/products/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", jsonContentType)
	return req
}

func newDeleteProductRequest(id int64) *http.Request {
	code := ProductCode{id}
	body, _ := json.Marshal(code)
	req, _ := http.NewRequest(http.MethodPost, "/products/delete", bytes.NewBuffer(body))
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

func assertDeals(t *testing.T, got, want []Deal) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertShoppingCart(t *testing.T, got, want ShoppingCart) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
