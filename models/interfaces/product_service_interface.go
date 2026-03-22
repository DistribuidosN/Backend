package interfaces

import "Backend/models"

// ProductServiceInterface defines the operations for products
type ProductServiceInterface interface {
	GetProducts() ([]models.Product, error)
	GetProduct(id int) (*models.Product, error)
	CreateProduct(product models.Product) (*models.Product, error)
	UpdateProduct(product models.Product) (*models.Product, error)
	DeleteProduct(id int) error
}
