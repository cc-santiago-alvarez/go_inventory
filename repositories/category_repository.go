package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := "INSERT INTO categories (name, prefix, description) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, category.Name, category.Prefix, category.Description).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	query := "SELECT * FROM categories ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Prefix, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) FindById(ctx context.Context, id string) (*models.Category, error) {
	query := "SELECT * FROM categories WHERE id = $1"

	var category models.Category
	err := r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.Prefix, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"

	result, err := r.db.ExecContext(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener el numero de filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("categoria no encontrada")
	}

	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM categories WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar la categoria: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener el numero de filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("categoria no encontrada")
	}

	return nil
}

func (r *CategoryRepository) PrefixExists(ctx context.Context, prefix string) (bool, error) {
    query := "SELECT EXISTS(SELECT 1 FROM categories WHERE prefix = $1)"
    var exists bool
    err := r.db.QueryRowContext(ctx, query, prefix).Scan(&exists)
    return exists, err
}
