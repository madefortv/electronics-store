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
	productService.repository.insertDeal(Deal{Name: "Half Off", Type: "Percent", Percent: "0.5"})
	productService.repository.insertDeal(Deal{Name: "Laptop Mouse Bundle", Type: "Bundle"})
	productService.repository.insertDeal(Deal{Name: "Buy 3 Get 2 free", Type: "BuyXGetY", X: 3, Y: 2})
	productService.repository.insertDeal(Deal{Name: "$10 keyboard", Type: "Coupon", Coupon: "10"})

	// some products to list
	productService.repository.createProductsTable()
	productService.repository.insertProduct(Product{1, "laptop", "very fast", "1000.00"})
	productService.repository.insertProduct(Product{2, "mouse", "much clicky", "10.00"})
	productService.repository.insertProduct(Product{3, "monitor", "four kay", "100.00"})
	productService.repository.insertProduct(Product{4, "usb", "type see", "5.00"})
	productService.repository.insertProduct(Product{5, "keyboard", "mecha", "25.00"})

	// actual items
	productService.repository.createOfferingsTable()
	productService.repository.insertOffering(Offering{ProductID: 2, DealID: 3, Active: true, ModifiedPrice: "1000.00"})
	productService.repository.insertOffering(Offering{ProductID: 1, DealID: 3, Active: true, ModifiedPrice: "1000.00"})
	productService.repository.insertOffering(Offering{ProductID: 3, DealID: 2, Active: true, ModifiedPrice: "NAN"})
	productService.repository.insertOffering(Offering{ProductID: 4, DealID: 4, Active: true, ModifiedPrice: "NAN"})
	productService.repository.insertOffering(Offering{ProductID: 5, DealID: 5, Active: true, ModifiedPrice: "NAN"})

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

		body, _ := json.Marshal(Product{ID: 3, Name: "monitor", Description: "four kay", Price: "100.00"})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 1}}
		want := ShoppingCart{Items: items, Total: "50"}

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

		body, _ := json.Marshal(Item{Product{3, "laptop", "very fast", "1000.00"}, 2})

		req, _ := http.NewRequest(http.MethodPut, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{3, "monitor", "four kay", "100.00"}, 2}}
		want := ShoppingCart{Items: items, Total: "100"}

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

	t.Run("Add an X of Y item to cart that doesn't meet threshold", func(t *testing.T) {

		body, _ := json.Marshal(Product{ID: 4})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 2},
			{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 1}}
		want := ShoppingCart{Items: items, Total: "105"}

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

	t.Run(" Trigger a buy x get y discount ", func(t *testing.T) {

		body, _ := json.Marshal(Item{Product{4, "useb", "type see", "5.00"}, 7})

		req, _ := http.NewRequest(http.MethodPut, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 2},
			{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 7}}

		want := ShoppingCart{Items: items, Total: "125"}
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

	t.Run("Add an item with coupon discount to cart", func(t *testing.T) {

		body, _ := json.Marshal(Product{ID: 5})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 2},
			{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 7},
			{Product{ID: 5, Name: "keyboard", Price: "25.00", Description: "mecha"}, 1}}
		want := ShoppingCart{Items: items, Total: "140"}

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

	t.Run("Add a bundle item to cart, there should be no effect", func(t *testing.T) {

		body, _ := json.Marshal(Product{ID: 1})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 2},
			{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 7},
			{Product{ID: 5, Name: "keyboard", Price: "25.00", Description: "mecha"}, 1},
			{Product{ID: 1, Name: "laptop", Price: "1000.00", Description: "very fast"}, 1}}
		want := ShoppingCart{Items: items, Total: "1140"}

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

	t.Run("Add the other bundled item to the cart", func(t *testing.T) {

		body, _ := json.Marshal(Product{ID: 2})

		req, _ := http.NewRequest(http.MethodPost, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 3, Name: "monitor", Price: "100.00", Description: "four kay"}, 2},
			{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 7},
			{Product{ID: 5, Name: "keyboard", Price: "25.00", Description: "mecha"}, 1},
			{Product{ID: 1, Name: "laptop", Price: "1000.00", Description: "very fast"}, 1},
			{Product{ID: 2, Name: "mouse", Price: "10.00", Description: "much clicky"}, 1}}
		want := ShoppingCart{Items: items, Total: "1140"}

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

	t.Run("Delete an item from cart", func(t *testing.T) {

		body, _ := json.Marshal(Product{ID: 3})

		req, _ := http.NewRequest(http.MethodDelete, "/cart", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		items := []Item{{Product{ID: 4, Name: "usb", Price: "5.00", Description: "type see"}, 7},
			{Product{ID: 5, Name: "keyboard", Price: "25.00", Description: "mecha"}, 1},
			{Product{ID: 1, Name: "laptop", Price: "1000.00", Description: "very fast"}, 1},
			{Product{ID: 2, Name: "mouse", Price: "10.00", Description: "much clicky"}, 1}}
		want := ShoppingCart{Items: items, Total: "1040"}

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
	productService.repository.insertOffering(Offering{ProductID: 2, DealID: 1})
	// laptop with a mouse free
	productService.repository.insertOffering(Offering{ProductID: 2, DealID: 3})
	productService.repository.insertOffering(Offering{ProductID: 1, DealID: 3})
	// 50% off monitors
	productService.repository.insertOffering(Offering{ProductID: 3, DealID: 2})

	t.Run("create new offering connecting usbs to the buy 3 USBs get 1 free offering", func(t *testing.T) {

		body, _ := json.Marshal(Offering{ProductID: 4, DealID: 4})
		req, _ := http.NewRequest(http.MethodPost, "/offerings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", jsonContentType)

		want := ""
		var got string
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, req)

		assertStatus(t, response.Code, http.StatusCreated)

		assertResponseBody(t, got, want)
	})

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

	t.Run("get the list of deals", func(t *testing.T) {

		request, _ := http.NewRequest(http.MethodGet, "/deals", nil)
		want := []Deal{{ID: 1, Name: "Regular Price", Type: "Retail"}, {ID: 2, Name: "Half Off", Type: "Percent", Percent: "50"}}

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

		request := newProductRequest(http.MethodPost, 0, "monitor", "fourkay", "100.00")
		want := ""
		var got string
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusCreated)

		assertResponseBody(t, got, want)

	})

	t.Run("update a product name and description", func(t *testing.T) {

		request := newProductRequest(http.MethodPut, 1, "laptop", "older", "85.00")
		want := ""

		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)

		var got string

		assertStatus(t, response.Code, http.StatusNoContent)
		assertResponseBody(t, got, want)

	})

	t.Run("delete the a product (id = 2)", func(t *testing.T) {
		request := newProductRequest(http.MethodDelete, 2, "monitor", "fourkay", "100.00")
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

func newProductRequest(method string, id int, name, description, price string) *http.Request {
	product := Product{
		id,
		name,
		description,
		price,
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(method, "/products", bytes.NewBuffer(body))
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
