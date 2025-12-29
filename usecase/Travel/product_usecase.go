package usecase

import (
	"errors"
	"multi-onboarding/model/Travel"
	"multi-onboarding/repository/Travel"
)

// ProductUsecase interface defines business logic operations
type ProductUsecase interface {
	// Draft operations
	AddDraft(draft *model.ProductDraft) error
	GetAllDrafts() ([]model.ProductDraft, error)
	ClearDrafts() error
	DeleteDraft(id int) error
	ConfirmDrafts(createdBy int64) error

	// Database operations
	GetAllProducts() ([]model.TravelProduct, error)
	GetProductByID(id int) (*model.TravelProduct, error)
	CreateProduct(product *model.TravelProduct) error
	UpdateProduct(id int, product *model.TravelProduct) error
	DeleteProduct(id int) error
}

// productUsecase implements ProductUsecase
type productUsecase struct {
	repo repository.ProductRepository
}

// NewProductUsecase creates a new product usecase
func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		repo: repo,
	}
}

// GetAllProducts retrieves all products
func (uc *productUsecase) GetAllProducts() ([]model.TravelProduct, error) {
	return uc.repo.GetAll()
}

// GetProductByID retrieves a product by ID
func (uc *productUsecase) GetProductByID(id int) (*model.TravelProduct, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	product, err := uc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

// CreateProduct creates a new product with validation
func (uc *productUsecase) CreateProduct(product *model.TravelProduct) error {
	if err := uc.validateProduct(product); err != nil {
		return err
	}
	return uc.repo.Create(product)
}

// UpdateProduct updates an existing product
func (uc *productUsecase) UpdateProduct(id int, product *model.TravelProduct) error {
	if id <= 0 {
		return ErrInvalidID
	}
	if err := uc.validateProduct(product); err != nil {
		return err
	}
	product.ID = int64(id)
	err := uc.repo.Update(product)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

// Draft operations

// AddDraft adds a draft product
func (uc *productUsecase) AddDraft(draft *model.ProductDraft) error {
	return uc.repo.AddDraft(draft)
}

// GetAllDrafts gets all draft products
func (uc *productUsecase) GetAllDrafts() ([]model.ProductDraft, error) {
	return uc.repo.GetAllDrafts()
}

// ClearDrafts clears all draft products
func (uc *productUsecase) ClearDrafts() error {
	return uc.repo.ClearDrafts()
}

// DeleteDraft deletes a draft product by ID
func (uc *productUsecase) DeleteDraft(id int) error {
	return uc.repo.DeleteDraft(id)
}

// ConfirmDrafts confirms all draft products to database
func (uc *productUsecase) ConfirmDrafts(createdBy int64) error {
	return uc.repo.ConfirmDrafts(createdBy)
}

// DeleteProduct deletes a product by ID
func (uc *productUsecase) DeleteProduct(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	err := uc.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

// validateProduct validates product data
func (uc *productUsecase) validateProduct(product *model.TravelProduct) error {
	if product.Code == "" {
		return ErrNameRequired // Reuse error for code
	}
	if product.Name == "" {
		return ErrNameRequired
	}
	if product.InsuranceCode == "" {
		return ErrNameRequired // Reuse error for insurance_code
	}
	if product.Logo == "" {
		return ErrNameRequired // Reuse error for logo
	}
	if product.Type == "" {
		return ErrNameRequired // Reuse error for type
	}
	return nil
}

// isValidStatus checks if status is valid
func isValidStatus(status string) bool {
	validStatuses := []string{"draft", "published", "archived"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

