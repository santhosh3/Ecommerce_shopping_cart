package product

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/santhosh3/ECOM/models"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAllProducts() ([]*models.Product, error) {
	var products []models.Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err // return the error instead of nil
	}

	// Convert `products` to a slice of pointers
	productPointers := make([]*models.Product, len(products))
	for i := range products {
		productPointers[i] = &products[i]
	}

	return productPointers, nil
}

func (s *Store) CreateProduct(product models.Product) (*models.Product, error) {
	if err := s.db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to create the product %v", err)
	}
	return &product, nil
}

func (s *Store) GetProductById(id int16) (*models.Product, error) {
	var product models.Product

	err := s.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found %s", err)
		}
		return nil, err
	}
	return &product, nil
}

func (s *Store) DeleteProductById(id int16) (*models.Product, error) {
	var product models.Product

	// First, find the product by ID
	err := s.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found: %v", err)
		}
		return nil, err
	}

	// Now, delete the product
	err = s.db.Delete(&product, id).Error
	if err != nil {
		return nil, err
	}

	folderPath := "./uploads/products"
	fileName := strings.Split(product.ProductImage, "/")[4]
	filePath := filepath.Join(folderPath, fileName)
	os.Remove(filePath)
	// Return the product that was deleted
	return &product, nil
}

func (s *Store) UpdateProductById(id uint64, productUpdate models.Product) (*models.Product, error) {
	var product models.Product
	folderPath := "./uploads/products"

	// Find the existing product by ID
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	if len(productUpdate.ProductImage) != 0 {
		//generate unique file name
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), productUpdate.ProductImage)
		filePath := filepath.Join(folderPath, fileName)
 
		productUpdate.ProductImage = fmt.Sprintf("/api/v1/productImage/%s", fileName)

		out, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("error While uploadinf file: %v", err)
		}
		defer out.Close()

		oldFileName := strings.Split(product.ProductImage, "/")[4]
		OldFilePath := filepath.Join(folderPath, oldFileName)

		os.Remove(OldFilePath)
	}

	// Update only the provided fields
	if err := s.db.Model(&product).Updates(productUpdate).Error; err != nil {
		return nil, err
	}

	// Return the updated product
	return &product, nil
}
