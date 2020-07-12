package main

import (
	"github.com/shopspring/decimal"
)

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

// use recursion to calculate the # of regular price items the charge for.
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

// Takes a map [DealID] -> ProductOffering{} to determine
func (service *ProductService) bundlePrice(bundledItems map[int][]*ProductOffering) (decimal.Decimal, error) {
	bundleTotal, err := decimal.NewFromString("0")
	if err != nil {
		return decimal.Decimal{}, err
	}
	// for each unique bundle in the cart
	for k, v := range bundledItems {
		// get all the items in that bundle
		offerings := service.repository.getBundleComponents(k)
		// start a subtotal for that bundle
		subtotal, err := decimal.NewFromString("0")
		if err != nil {
			return decimal.Decimal{}, err
		}
		for i := range v {
			if len(offerings) == len(bundledItems[k]) {
				price, err := decimal.NewFromString(v[i].ModifiedPrice)
				if err != nil {
					return decimal.Decimal{}, err
				}
				subtotal = subtotal.Add(price)
				break

			} else {
				/* add up the retail prices x quantity */
				price, err := decimal.NewFromString(v[i].Price)
				if err != nil {
					return decimal.Decimal{}, err
				}
				subtotal = subtotal.Add(price)
			}

		}

		bundleTotal = bundleTotal.Add(subtotal)
	}

	return bundleTotal, nil
}

func (service *ProductService) totalPrice(productOfferings []*ProductOffering) (string, error) {
	total, err := decimal.NewFromString("0")
	if err != nil {
		return "NAN", err
	}
	/* A map of deals to cart items that match those deals, used to calculate bundle prices */
	bundledItems := make(map[int][]*ProductOffering)
	/* a bit hacky, but this map tracks the count of each product in a bundle*/
	for i := range productOfferings {
		temp, err := decimal.NewFromString("0")
		if err != nil {
			return "NAN", err
		}
		po := *productOfferings[i]

		switch po.Type {
		case "Bundle":
			// we don't decided on a final price until we finish looping over all the products
			// so we add the each item in the cart to a list associated with the bundle it's in.
			bundledItems[po.DealID] = append(bundledItems[po.DealID], &po)
		case "BuyXGetY":

			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			// recurse over the number of items to calculate full price items
			regularPriceItems := buyXGetYPrice(po.Quantity, po.X, po.Y)
			quantity := decimal.NewFromInt(int64(regularPriceItems))
			temp = price.Mul(quantity)

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

	// calcualte the bundle price at the end and tac it on
	bundledTotal, err := service.bundlePrice(bundledItems)
	if err != nil {
		return "NAN", err
	}
	total = total.Add(bundledTotal)
	return total.String(), nil
}
