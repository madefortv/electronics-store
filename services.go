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

func (service *ProductService) FindAll() []*Product {
	if service.config.Enabled {
		return service.repository.FindAll()
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
