package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
	"github.com/cc-santiago-alvarez/go_inventory.git/repositories"
)

type CategoryService struct {
	repo        *repositories.CategoryRepository
	productRepo *repositories.ProductRepository
}

func NewCategoryService(repo *repositories.CategoryRepository, productRepo *repositories.ProductRepository) *CategoryService {
	return &CategoryService{repo: repo, productRepo: productRepo}
}

func (c *CategoryService) FindCategoryWithProducts(ctx context.Context, id string) (*models.CategoryWithProducts, error) {
	category, err := c.repo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la categoría: %w", err)
	}

	products, err := c.productRepo.FindByCategoryID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los productos de la categoría: %w", err)
	}

	return &models.CategoryWithProducts{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		Products:    products,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (c *CategoryService) CreateCategory(ctx context.Context, name, description string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Generar prefijo único
	prefix, err := c.resolveUniquePrefix(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error generando prefijo: %w", err)
	}

	category := &models.Category{
		Name:        name,
		Prefix:      prefix,
		Description: description,
	}

	if err := c.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (c *CategoryService) FindAllCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := c.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las categorias: %w", err)
	}

	return categories, nil
}

func (c *CategoryService) FindCategoryById(ctx context.Context, id string) (*models.Category, error) {
	category, err := c.repo.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la categoria: %w", err)
	}

	return category, nil
}

func (c *CategoryService) UpdateCategory(ctx context.Context, id string, name, description string) error {
	if name == "" {
		return errors.New("El nombre de la categoria es requerido")
	}

	category, err := c.repo.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("Error al obtener la categoria: %w", err)
	}

	category.Name = name
	category.Description = description

	err = c.repo.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("Error al actualizar la categoria: %w", err)
	}

	return nil
}

func (c *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	_, err := c.repo.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("Error al obtener la categoria: %w", err)
	}

	err = c.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("Error al eliminar la categoria: %w", err)
	}

	return nil
}

func generatePrefix(name string) string {
	// Convertir a mayúsculas y tomar los primeros 3 caracteres
	normalized := strings.ToUpper(strings.TrimSpace(name))

	// Remover caracteres no alfanuméricos (tildes, espacios, etc.)
	var clean []rune
	for _, r := range normalized {
		if r >= 'A' && r <= 'Z' {
			clean = append(clean, r)
		}
	}

	if len(clean) >= 3 {
		return string(clean[:3])
	}
	return string(clean)
}

func (c *CategoryService) resolveUniquePrefix(ctx context.Context, name string) (string, error) {
	base := generatePrefix(name)
	candidate := base

	for i := 1; i <= 99; i++ {
		exists, err := c.repo.PrefixExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		// ELE → EL1 → EL2 → ... → E10 → E11 ...
		suffix := fmt.Sprintf("%d", i)
		candidate = base[:max(1, len(base)-len(suffix))] + suffix
	}

	return "", errors.New("no se pudo generar un prefijo único")
}
