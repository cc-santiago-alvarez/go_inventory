package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
	"github.com/cc-santiago-alvarez/go_inventory.git/repositories"
)

type ProductService struct {
	repo         *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

func (p *ProductService) CreateProduct(ctx context.Context, name, description string, price float64, quantity int, categoryID string) (*models.Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	//if categoryID == "" {
	//	return nil, errors.New("category_id es requerido")
	//}

	//_, err := p.categoryRepo.FindById(ctx, categoryID)
	//if err != nil {
	//	return nil, fmt.Errorf("la categoría no existe: %w", err)
	//}

	category, err := p.categoryRepo.FindById(ctx, categoryID)
    if err != nil {
        return nil, fmt.Errorf("la categoría no existe: %w", err)
    }

	code, err := p.generateProductCode(ctx, categoryID)
    if err != nil {
        return nil, fmt.Errorf("error generando código: %w", err)
    }

	existing, err := p.repo.FindByNameAndCategory(ctx, name, categoryID)
    if err == nil && existing != nil {
        existing.Quantity += quantity
        if err := p.repo.Update(ctx, existing); err != nil {
            return nil, fmt.Errorf("error actualizando cantidad: %w", err)
        }
        return existing, nil
    }

	product := &models.Product{
        Name:        name,
        Code:        code,
        Description: description,
        Price:       price,
        Quantity:    quantity,
        Category: models.ProductCategory{
            ID:   category.ID,
            Name: category.Name,
            Prefix: category.Prefix,
        },
    }

	if err := p.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil

}

func (p *ProductService) FindAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := p.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los productos: %w", err)
	}

	return products, nil
}

func (p *ProductService) FindProductById(ctx context.Context, productID string) (*models.Product, error) {
	product, err := p.repo.FindById(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el producto: %w", err)
	}

	return product, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, productID string, name, description string, price float64, quantity int, categoryID string) error {
	if name == "" {
		return errors.New("El nombre del producto es requerido")
	}

	if price <= 0 {
		return errors.New("El precio del producto debe ser mayor a 0")
	}

	product, err := p.repo.FindById(ctx, productID)
	if err != nil {
		return fmt.Errorf("Error al obtener el producto: %w", err)
	}

	if categoryID == "" {
		return errors.New("category_id es requerido")
	}

	_, err = p.categoryRepo.FindById(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("la categoría no existe: %w", err)
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.Quantity = quantity
	product.Category.ID = categoryID

	err = p.repo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("Error al actualizar el producto: %w", err)
	}

	return nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, productID string) error {
	_, err := p.repo.FindById(ctx, productID)
	if err != nil {
		return fmt.Errorf("Error al obtener el producto: %w", err)
	}

	if err := p.repo.Delete(ctx, productID); err != nil {
		return fmt.Errorf("Error al eliminar el producto: %w", err)
	}

	return nil
}

func (p *ProductService) generateProductCode(ctx context.Context, categoryID string) (string, error) {
    // 1. Obtener el prefijo de la categoría
    category, err := p.categoryRepo.FindById(ctx, categoryID)
    if err != nil {
        return "", fmt.Errorf("categoría no encontrada: %w", err)
    }

    // 2. Obtener el número más alto existente para esa categoría
    maxCode, err := p.repo.MaxCodeByCategory(ctx, categoryID)
    if err != nil {
        return "", err
    }

    // 3. Generar el código: PREFIX-0001
    code := fmt.Sprintf("%s-%04d", category.Prefix, maxCode+1)
    return code, nil
}
