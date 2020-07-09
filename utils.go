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

func fitlerOfferByDealType(productOfferings []*ProductOffering, dtype DealType) []*ProductOffering {

	var offers []*ProductOffering
	for i := range productOfferings {
		po := *productOfferings[i]
		// might need deep reflect
		if po.Type == dtype {
			offers = append(offers, &po)
		}
	}
	return offers
}

// returns the # of regular price items
func buyXGetYPrice(quantity int, x int, y int) int {
	z := x + y
	// just in case
	if quantity < 0 {
		return 0
	}
	// base case: no deal applies to these items
	if quantity < x {
		return quantity
	}
	// base case: pay just for x of y items
	if quantity <= z {
		return x
	}
	// if we have more items, see if we can have deals
	return x + buyXGetYPrice(quantity-z, x, y)
}
