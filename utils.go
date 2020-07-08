package main

import (
	"log"
	"reflect"
)

func itemIndexByProduct(items []Item, product Product) int {

	log.Printf("comparing %v and %v", items, product)
	for i := range items {
		// might need deep reflect
		if reflect.DeepEqual(&items[i].Product, &product) {
			return i
		}
	}
	return -1
}
