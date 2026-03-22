package services

import (
	"Backend/models"
	"Backend/models/interfaces"
	"net/http"
	"sync"
)

var _ interfaces.ProductServiceInterface = (*ProductService)(nil)

// ProductService is the implementation of ProductServiceInterface
type ProductService struct {
	products map[int]models.Product
	nextID   int
	mu       sync.Mutex
}

// NewProductService creates a new ProductService
func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[int]models.Product),
		nextID:   1,
	}
}

// GetProducts returns all products
func (s *ProductService) GetProducts() ([]models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var products []models.Product
	for _, p := range s.products {
		products = append(products, p)
	}
	return products, nil
}

// GetProduct returns a product by ID
func (s *ProductService) GetProduct(id int) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.products[id]
	if !ok {
		return nil, models.NewApiError(http.StatusNotFound, "product not found")
	}
	return &product, nil
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product models.Product) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	product.ID = s.nextID
	s.products[product.ID] = product
	s.nextID++
	return &product, nil
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(product models.Product) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.products[product.ID]
	if !ok {
		return nil, models.NewApiError(http.StatusNotFound, "product not found")
	}

	s.products[product.ID] = product
	return &product, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.products[id]
	if !ok {
		return models.NewApiError(http.StatusNotFound, "product not found")
	}

	delete(s.products, id)
	return nil
}
