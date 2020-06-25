package main

import (
	"fmt"
	"net/http"
)

func StoreServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
