package main

/*
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

		assertShoppingCart(t, shoppingCart, want)
	})

	t.Run("Create a new shopping cart", func(t *testing.T) {

		items := make([]Item, 0)
		newCart := ShoppingCart{Items: items, Total: "0.0"}

		bytes, err := json.Marshal(newCart)
		if err != nil {
			t.Fatalf("Unable to marshal shopping cart")
		}

		var want ShoppingCart
		json.Unmarshal(bytes, want)

		assertShoppingCart(t, newCart, want)
	})
}

func UtilsTest(t *testing.T) {
	t.Run("Can find product from Intems list", func(t *testing.T) {
		items := []Item{{Product{1, "laptop", "very fast", "1000.00"}, 1}, {Product{2, "mouse", "clicky", "10.00"}, 1}}
		index := itemIndexByProduct(items, Product{1, "laptop", "very fast", "1000.00"})
		if index != 0 {
			t.Error("Got the wrong index")
		}

		notFound := itemIndexByProduct(items, Product{4, "DNE", "Eww", "1000.00"})
		if notFound != -1 {
			t.Error("Got wront index")
		}

	})
}
*/
