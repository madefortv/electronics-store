package main

import (
	"errors"
	"github.com/shopspring/decimal"
	"log"
)

type ProductService struct {
	config     *Config
	repository *ProductRepository
}

func NewProductService(config *Config, repository *ProductRepository) *ProductService {
	return &ProductService{config: config, repository: repository}
}

/* Shopping Cart */
func (service *ProductService) listCartItems() []Item {
	return service.repository.listCart()
}

func (service *ProductService) addToCart(product Product) error {
	return service.repository.addToCart(product)
}

func (service *ProductService) updateCart(item Item) error {
	return service.repository.updateCart(item)
}

func (service *ProductService) calculateTotalPrice() (string, error) {

	productOfferings := service.repository.getProductOfferings()

	total, err := decimal.NewFromString("0")
	if err != nil {
		return "NAN", err
	}

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
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			bundlePrice, err := decimal.NewFromString(po.ModifiedPrice)
			if err != nil {
				return "NAN", err
			}

			//items := service.repository.listCart()
			if price.LessThan(bundlePrice) {
				// this isn't the big ticket item
			}
			temp = price
			log.Printf("total: %s", temp.String())

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

			percent, err := decimal.NewFromString(po.Percent)
			if err != nil {
				return "NAN", err
			}
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			quantity := decimal.NewFromInt(int64(po.Quantity))
			//Note: percent should be in the range (0, 1)
			//TODO: add validation
			temp = price.Mul(percent).Mul(quantity)
			log.Printf("total: %s", temp.String())

		case "Coupon":

			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			coupon, err := decimal.NewFromString(po.Coupon)
			if err != nil {
				return "NAN", err
			}

			temp = price.Sub(coupon)

		default:
			price, err := decimal.NewFromString(po.Price)
			if err != nil {
				return "NAN", err
			}
			temp = price

		}

		total = total.Add(temp)
		log.Printf("total: %s", total.String())
	}
	return total.String(), nil
}

/* Products */
func (service *ProductService) getProduct(product Product) Product {
	if service.config.Enabled {
		// check if the product exists, otherwise return empty
		product, err := service.repository.getProduct(product)
		if err != nil {
			return Product{}
		}
		return product
	}
	return Product{}

}

func (service *ProductService) listProducts() []*Product {
	if service.config.Enabled {
		return service.repository.listProducts()
	}
	return []*Product{}
}

func (service *ProductService) newProduct(product Product) error {
	if service.config.Enabled {
		return service.repository.insertProduct(product)
	}
	return errors.New("Operation Not Permitted")

}

func (service *ProductService) updateProduct(product Product) error {
	if service.config.Enabled {
		return service.repository.updateProduct(product)
	}
	return errors.New("Operation Not Permitted")
}

func (service *ProductService) deleteProduct(code ProductCode) error {
	if service.config.Enabled {
		return service.repository.deleteProduct(code)
	}
	return errors.New("Operation Not Permitted")
}

/* Deals */
func (service *ProductService) newDeal(deal Deal) error {
	if service.config.Enabled {
		return service.repository.insertDeal(deal)
	}
	return errors.New("Operation Not Permitted")
}

func (service *ProductService) listDeals() []*Deal {
	if service.config.Enabled {
		return service.repository.listDeals()
	}
	return []*Deal{}
}

/* Offerings */
func (service *ProductService) newOffering(offering Offering) error {
	if service.config.Enabled {
		return service.repository.insertOffering(offering)
	}
	return errors.New("Operation Not Permitted")
}
