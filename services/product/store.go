package product

import (
	"context"
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

func (s *Store) withTimeout(query func(db *gorm.DB) error) error {
	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Apply the context to the DB connection
	dbWithTimeout := s.db.WithContext(ctx)

	// Execute the query function passed as argument
	err := query(dbWithTimeout)

	// Check for deadline exceeded error
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("not able to fetch from DB: %w", err)
	}

	return err
}

func (s *Store) GetFilteredProducts(size, name, priceGreaterThan, priceLessThan, priceSort string) ([]*models.Product, error) {
	var products []models.Product

	// Create a context with a 1 second timeout for DB queries
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Apply the context with timeout to the DB connection
	dbWithTimeout := s.db.WithContext(ctx)

	// Start building the query with a base filter for non-deleted products
	filterQuery := dbWithTimeout.Where("is_deleted = ?", false)

	// Filter by size
	if size != "" {
		filterQuery = filterQuery.Where("available_sizes = ?", size)
	}
	// Filter by name
	if name != "" {
		filterQuery = filterQuery.Where("title LIKE ?", "%"+name+"%")
	}

	// Filter by price range
	if priceGreaterThan != "" && priceLessThan != "" {
		filterQuery = filterQuery.Where("price BETWEEN ? AND ?", priceGreaterThan, priceLessThan)
	} else if priceGreaterThan != "" {
		filterQuery = filterQuery.Where("price >= ?", priceGreaterThan)
	} else if priceLessThan != "" {
		filterQuery = filterQuery.Where("price <= ?", priceLessThan)
	}

	// Sort by price
	if priceSort != "" {
		sortOrder := "price ASC"
		if priceSort == "-1" {
			sortOrder = "price DESC"
		}
		filterQuery = filterQuery.Order(sortOrder)
	}

	// Execute query and return results
	if err := filterQuery.Find(&products).Error; err != nil {
		return nil, err
	}

	// Convert `products` to a slice of pointers
	productPointers := make([]*models.Product, len(products))
	for i := range products {
		productPointers[i] = &products[i]
	}

	return productPointers, nil
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
     err := s.withTimeout(func(db *gorm.DB) error {
		return db.Create(&product).Error;
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to create the product %s", err)
		}
		return nil, err
	}

	return &product, nil
}

func (s *Store) GetProductById(id int16) (*models.Product, error) {
	var product models.Product

	// Use the helper to execute the query with a 1-second timeout
	err := s.withTimeout(func(db *gorm.DB) error {
		return db.First(&product, id).Error
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
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
	err := s.withTimeout(func(db *gorm.DB) error {
		return db.First(&product, id).Error;
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%v", err)
		}
		return nil, err
	}

	// Now, delete the product
	err = s.withTimeout(func(db *gorm.DB) error {
		return db.Delete(&product, id).Error
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
		return nil, err
	}

	go func() {
		folderPath := "./uploads/products"
		fileName := strings.Split(product.ProductImage, "/")[4]
		filePath := filepath.Join(folderPath, fileName)
		os.Remove(filePath)
	}()

	// Return the product that was deleted
	return &product, nil
}

func (s *Store) UpdateProductById(id uint64, productUpdate models.Product) (*models.Product, error) {
	var product models.Product
	folderPath := "./uploads/products"

	// Find the existing product by ID
	if err := s.withTimeout(func(db *gorm.DB) error {
		return db.First(&product, id).Error;
	}); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
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
	err := s.withTimeout(func(db *gorm.DB) error {
		return s.db.Model(&product).Updates(productUpdate).Error
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%s", err)
		}
		return nil, err
	}
	
	// Return the updated product
	return &product, nil
}
