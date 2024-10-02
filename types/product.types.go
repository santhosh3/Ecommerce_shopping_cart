package types

import "github.com/santhosh3/ECOM/models"


type ProductStore interface {
	GetAllProducts() ([]*models.Product, error)
	CreateProduct(product models.Product) (*models.Product, error)
	GetProductById(id int16) (*models.Product, error)
	DeleteProductById(id int16) (*models.Product, error)
	UpdateProductById(id uint64, productUpdate models.Product) (*models.Product, error)
	GetFilteredProducts(size, name, priceGreaterThan, priceLessThan, priceSort string) ([]*models.Product, error)
}
