package main

import (
	"log"

	"github.com/shopspring/decimal"
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

func filterOfferByDealId(productOfferings []*ProductOffering, dealId int) []*ProductOffering {

	var offers []*ProductOffering
	for i := range productOfferings {
		po := *productOfferings[i]
		// might need deep reflect
		if po.DealId == dealId {
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

func totalPrice(productOfferings []*ProductOffering) (string, error) {
	total, err := decimal.NewFromString("0")
	if err != nil {
		return "NAN", err
	}
	retailPriceMap := make(map[int]decimal.Decimal)
	bundlePriceMap := make(map[int]decimal.Decimal)
	bestPriceMap := make(map[int]decimal.Decimal)
	for i := range productOfferings {
		temp, err := decimal.NewFromString("0")
		if err != nil {
			return "NAN", err
		}

		if err != nil {
			return "NAN", err
		}
		//var n int // the number of an item to reduce
		po := *productOfferings[i]
		switch po.Type {
		case "Bundle":
			// basically just sum up each deal/product
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			modifiedPrice, ok := retailPriceMap[po.DealId]
			quantity := decimal.NewFromInt(int64(po.Quantity))
			temp := price.Mul(quantity)
			if !ok {
				retailPriceMap[po.DealId] = temp
			} else {
				retailPriceMap[po.DealId] = modifiedPrice.Add(temp)
			}
			// we don't decided on a final price until we finish looping over all the products
		case "BuyXGetY":

			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			// recurse over the number of items to calculate full price items
			regularPriceItems := buyXGetYPrice(po.Quantity, po.X, po.Y)
			quantity := decimal.NewFromInt(int64(regularPriceItems))
			temp = price.Mul(quantity)
			log.Printf("total: %s", temp.String())
			//remove

		case "Percent":

			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			// should be in range (0,1)
			percent, err := decimal.NewFromString(po.Percent)
			if err != nil {
				return "NAN", err
			}
			quantity := decimal.NewFromInt(int64(po.Quantity))
			temp = price.Mul(quantity).Mul(percent)

		case "Coupon":
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			// should be in range (0,1)
			coupon, err := decimal.NewFromString(po.Coupon)
			if err != nil {
				return "NAN", err
			}
			quantity := decimal.NewFromInt(int64(po.Quantity))
			discountedPrice := price.Sub(coupon)
			temp = discountedPrice.Mul(quantity)

		case "Retail":
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			quantity := decimal.NewFromInt(int64(po.Quantity))
			temp = price.Mul(quantity)

		}

		total = total.Add(temp)
	}

	offers := fitlerOfferByDealType(productOfferings, "Bundle")

	// loop over the bundles, compare retailPrice to bundledPrice
	for i := range offers {
		offer := *offers[i]
		modifiedPrice, ok := retailPriceMap[offer.DealId]
		if !ok {
			return "NAN", err
		}

		quantity := decimal.NewFromInt(int64(offer.Quantity))

		bundlePrice, err := decimal.NewFromString(offer.ModifiedPrice)
		totalBundlePrice := bundlePrice.Mul(quantity)
		if err != nil {
			return "NAN", err
		}

		bundlePriceMap[offer.DealId] = bundlePrice.Mul(quantity)

		bestPriceMap[offer.DealId] = decimal.Min(totalBundlePrice, modifiedPrice)

	}

	// add the best prices for the bundles
	for _, v := range bestPriceMap {
		total = total.Add(v)
	}

	return total.String(), nil
}
