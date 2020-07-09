package main

import (
	"errors"
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
	total, err := totalPrice(productOfferings)
	if err != nil {
		return "NAN", err
	}
	return total, nil
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
