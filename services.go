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

/* Products */
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

/*
func (service *ProductService) listOfferings() []*Offering {
	if service.config.Enabled {
		return service.repository.listOfferings()
	}
	return []*Offering{}
}

*/
