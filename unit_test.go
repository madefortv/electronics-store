package main

import (
	"encoding/json"
	"testing"
)

/*

type ShoppingCart struct {
	Items []Item `json:"items"`
}

type Item struct {
	Product Product `json:"product"`
	Count   int     `json:"count"`
}

*/
func ShoppingCartModelTests(t *testing.T) {
	t.Run("test marshalling and unmarshalling shopping carts", func(t *testing.T) {
		shoppingCartItems := []byte(`[{"product": {"name": "laptop", "description": "very fast", "price": "1000.00"}, "count": 1},
					 {"product": {"name": "mouse", "description": "clicky ", "price": "10.00"}, "count": 1}]`)

		items := make([]Item, 0)
		json.Unmarshal(shoppingCartItems, items)

		shoppingCart := ShoppingCart{Items: items}

		bytes, err := json.Marshal(shoppingCart)
		if err != nil {
			t.Fatalf("Unable to marshal shopping cart")
		}

		var want ShoppingCart

		json.Unmarshal(bytes, want)
	})

}
