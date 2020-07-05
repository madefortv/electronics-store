package store

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

func (service *ProductService) createDeal(deal Deal) error {
	if service.config.Enabled {
		return service.repository.insertDeal(deal)
	}
	return errors.New("Operation not permitted")
}

func (service *ProductService) listProducts() []*Product {
	if service.config.Enabled {
		return service.repository.FindAll()
	}
	return []*Product{}
}

func (service *ProductService) createProduct(product Product) error {
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

func (service *ProductService) deleteProduct(product Product) error {
	if service.config.Enabled {
		return service.repository.deleteProduct(product)
	}
	return errors.New("Operation Not Permitted")
}
