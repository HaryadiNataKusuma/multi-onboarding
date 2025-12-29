package repository

import (
	"sync"
	"time"

	"multi-onboarding/model/Travel"
)

// ProductRepository interface defines the contract for product data operations
type ProductRepository interface {
	// Draft operations
	AddDraft(draft *model.ProductDraft) error
	GetAllDrafts() ([]model.ProductDraft, error)
	ClearDrafts() error
	DeleteDraft(id int) error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAll() ([]model.TravelProduct, error)
	GetByID(id int) (*model.TravelProduct, error)
	Create(product *model.TravelProduct) error
	Update(product *model.TravelProduct) error
	Delete(id int) error
}

// InMemoryProductRepository implements ProductRepository using in-memory storage
type InMemoryProductRepository struct {
	products []model.TravelProduct
	mu       sync.RWMutex
	nextID   int
}

// NewInMemoryProductRepository creates a new in-memory product repository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make([]model.TravelProduct, 0),
		nextID:   1,
	}
}

// GetAll retrieves all products
func (r *InMemoryProductRepository) GetAll() ([]model.TravelProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modifications
	result := make([]model.TravelProduct, len(r.products))
	copy(result, r.products)
	return result, nil
}

// GetByID retrieves a product by ID
func (r *InMemoryProductRepository) GetByID(id int) (*model.TravelProduct, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.products {
		if r.products[i].ID == int64(id) {
			// Return a copy
			product := r.products[i]
			return &product, nil
		}
	}
	return nil, ErrProductNotFound
}

// Create creates a new product
func (r *InMemoryProductRepository) Create(product *model.TravelProduct) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product.ID = int64(r.nextID)
	r.nextID++
	product.CreatedAt = time.Now()
	now := time.Now()
	product.UpdatedAt = &now
	// Status field removed from model

	r.products = append(r.products, *product)
	return nil
}

// Update updates an existing product
func (r *InMemoryProductRepository) Update(product *model.TravelProduct) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.products {
		if r.products[i].ID == product.ID {
			now := time.Now()
			product.UpdatedAt = &now
			// Preserve CreatedAt
			product.CreatedAt = r.products[i].CreatedAt
			r.products[i] = *product
			return nil
		}
	}
	return ErrProductNotFound
}

// Delete deletes a product by ID
func (r *InMemoryProductRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, product := range r.products {
		if product.ID == int64(id) {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return ErrProductNotFound
}

