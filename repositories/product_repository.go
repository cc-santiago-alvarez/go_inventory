package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := "INSERT INTO products (code, name, description, price, quantity, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, product.Code, product.Name, product.Description, product.Price, product.Quantity, product.Category.ID).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

// func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
// 	query := "SELECT * FROM products ORDER BY created_at DESC"

// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	products := []models.Product{}
// 	for rows.Next() {
// 		var product models.Product
// 		err := rows.Scan(&product.ID, &product.Code, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		products = append(products, product)
// 	}

// 	return products, nil
// }

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
    query := `
        SELECT p.id, p.code, p.name, p.description, p.price, p.quantity,
               c.id, c.name, c.prefix,
               p.created_at, p.updated_at
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        ORDER BY p.created_at DESC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    products := []models.Product{}
    for rows.Next() {
        var product models.Product
        err := rows.Scan(
            &product.ID, &product.Code, &product.Name, &product.Description,
            &product.Price, &product.Quantity,
            &product.Category.ID, &product.Category.Name, &product.Category.Prefix,
            &product.CreatedAt, &product.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

    return products, nil
}

func (r *ProductRepository) FindById(ctx context.Context, id string) (*models.Product, error) {
	query := "SELECT * FROM products WHERE id = $1"

	var product models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Code, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.Category.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	query := "UPDATE products SET name = $1, description = $2, price = $3, quantity = $4, category_id = $5 WHERE id = $6"

	result, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Quantity, product.Category.ID, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener el numero de filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pructo no encontrado")
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar el post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener el numero de filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pructo no encontrado")
	}

	return nil

}

func (r *ProductRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]models.Product, error) {
	query := "SELECT id, code, name, description, price, quantity, category_id, created_at, updated_at FROM products WHERE category_id = $1 ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Code, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.Category.ID, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) CountByCategory(ctx context.Context, categoryID string) (int, error) {
    query := "SELECT COUNT(*) FROM products WHERE category_id = $1"
    var count int
    err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&count)
    return count, err
}

func (r *ProductRepository) MaxCodeByCategory(ctx context.Context, categoryID string) (int, error) {
    query := "SELECT COALESCE(MAX(CAST(SPLIT_PART(code, '-', 2) AS INTEGER)), 0) FROM products WHERE category_id = $1"
    var maxCode int
    err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&maxCode)
    return maxCode, err
}

func (r *ProductRepository) FindByNameAndCategory(ctx context.Context, name, categoryID string) (*models.Product, error) {
    query := "SELECT id, code, name, description, price, quantity, category_id, created_at, updated_at FROM products WHERE name = $1 AND category_id = $2"
    var product models.Product
    err := r.db.QueryRowContext(ctx, query, name, categoryID).Scan(
        &product.ID, &product.Code, &product.Name, &product.Description,
        &product.Price, &product.Quantity, &product.Category.ID,
        &product.CreatedAt, &product.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &product, nil
}