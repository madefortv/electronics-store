package store

/*
import (
	"testing"

	"database/sql"
)

type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type MockDB struct{}

func (mdb *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	mdb.callParams = []interface{}{query}
	mdb.callparams = append(mdb.callParams, args...)

	return nil, nil
}

func (mdb *MockDB) CalledWith() []interface{} {
	return mdb.callParams
}

func TestDatabase(t *testing.T) {
	config := NewConfig()
	db, err := ConnectDatabase(config)
	if err != nil {
		panic(err)
	}
	productRepository := NewProductRepository(db)
	productService := NewProductService(config, productRepository)
	server := NewServer(config, productService)

	t.Run("inserts a new product into db", func(t *testing.T) {

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

}
*/
