package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETProduct(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	StoreServer(response, request)

	t.Run("returns product id", func(t *testing.T) {
		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
