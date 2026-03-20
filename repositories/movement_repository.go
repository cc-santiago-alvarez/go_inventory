package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
)

type MovementRepository struct {
	db *sql.DB
}

func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

func (r *MovementRepository) Create(ctx context.Context, movement *models.Movement) error {
	query := "INSERT INTO movements (product_id, type, quantity, reason, date) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, movement.Product.ID, movement.Type, movement.Quantity, movement.Reason, movement.Date).Scan(&movement.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MovementRepository) FindAll(ctx context.Context) ([]models.Movement, error) {
	query := `
		SELECT m.id, m.product_id, p.name, p.description, c.prefix,
		       m.type, m.quantity, m.reason, m.date, m.created_at
		FROM movements m
		LEFT JOIN products p ON m.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.Movement{}
	for rows.Next() {
		var movement models.Movement
		var productID, productName, productDesc, productPrefix sql.NullString
		err := rows.Scan(&movement.ID, &productID, &productName, &productDesc, &productPrefix,
			&movement.Type, &movement.Quantity, &movement.Reason, &movement.Date, &movement.CreatedAt)
		if err != nil {
			return nil, err
		}
		if productID.Valid {
			movement.Product.ID = productID.String
		}
		if productName.Valid {
			movement.Product.Name = productName.String
		}
		if productDesc.Valid {
			movement.Product.Description = productDesc.String
		}
		if productPrefix.Valid {
			movement.Product.Prefix = productPrefix.String
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func (r *MovementRepository) FindByProductID(ctx context.Context, productID string) ([]models.Movement, error) {
	query := `
		SELECT m.id, m.product_id, p.name, p.description, c.prefix,
		       m.type, m.quantity, m.reason, m.date, m.created_at
		FROM movements m
		LEFT JOIN products p ON m.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE m.product_id = $1
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.Movement{}
	for rows.Next() {
		var movement models.Movement
		var pID, productName, productDesc, productPrefix sql.NullString
		err := rows.Scan(&movement.ID, &pID, &productName, &productDesc, &productPrefix,
			&movement.Type, &movement.Quantity, &movement.Reason, &movement.Date, &movement.CreatedAt)
		if err != nil {
			return nil, err
		}
		if pID.Valid {
			movement.Product.ID = pID.String
		}
		if productName.Valid {
			movement.Product.Name = productName.String
		}
		if productDesc.Valid {
			movement.Product.Description = productDesc.String
		}
		if productPrefix.Valid {
			movement.Product.Prefix = productPrefix.String
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func (r *MovementRepository) FindByType(ctx context.Context, movementType models.MovementType) ([]models.Movement, error) {
	query := `
		SELECT m.id, m.product_id, p.name, p.description, c.prefix,
		       m.type, m.quantity, m.reason, m.date, m.created_at
		FROM movements m
		LEFT JOIN products p ON m.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE m.type = $1
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, movementType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.Movement{}
	for rows.Next() {
		var movement models.Movement
		var productID, productName, productDesc, productPrefix sql.NullString
		err := rows.Scan(&movement.ID, &productID, &productName, &productDesc, &productPrefix,
			&movement.Type, &movement.Quantity, &movement.Reason, &movement.Date, &movement.CreatedAt)
		if err != nil {
			return nil, err
		}
		if productID.Valid {
			movement.Product.ID = productID.String
		}
		if productName.Valid {
			movement.Product.Name = productName.String
		}
		if productDesc.Valid {
			movement.Product.Description = productDesc.String
		}
		if productPrefix.Valid {
			movement.Product.Prefix = productPrefix.String
		}
		movements = append(movements, movement)
	}
	return movements, nil
}

func (r *MovementRepository) FindByDateRange(ctx context.Context, from time.Time, to time.Time) ([]models.Movement, error) {
	query := `
		SELECT m.id, m.product_id, p.name, p.description, c.prefix,
		       m.type, m.quantity, m.reason, m.date, m.created_at
		FROM movements m
		LEFT JOIN products p ON m.product_id = p.id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE m.date BETWEEN $1 AND $2
		ORDER BY m.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := []models.Movement{}
	for rows.Next() {
		var movement models.Movement
		var productID, productName, productDesc, productPrefix sql.NullString
		err := rows.Scan(&movement.ID, &productID, &productName, &productDesc, &productPrefix,
			&movement.Type, &movement.Quantity, &movement.Reason, &movement.Date, &movement.CreatedAt)
		if err != nil {
			return nil, err
		}
		if productID.Valid {
			movement.Product.ID = productID.String
		}
		if productName.Valid {
			movement.Product.Name = productName.String
		}
		if productDesc.Valid {
			movement.Product.Description = productDesc.String
		}
		if productPrefix.Valid {
			movement.Product.Prefix = productPrefix.String
		}
		movements = append(movements, movement)
	}
	return movements, nil
}
