package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
	"github.com/cc-santiago-alvarez/go_inventory.git/repositories"
)

type MovementService struct {
	movementRepo *repositories.MovementRepository
	productRepo  *repositories.ProductRepository
}

func NewMovementService(repo *repositories.MovementRepository, productRepo *repositories.ProductRepository) *MovementService {
	return &MovementService{movementRepo: repo, productRepo: productRepo}
}

func (s *MovementService) CreateMovement(ctx context.Context, req models.CreateMovementRequest) (*models.Movement, error) {
	if req.ProductID == "" {
		return nil, errors.New("product_id es requerido")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a 0")
	}
	if req.Reason == "" {
		return nil, errors.New("el motivo es requerido")
	}
	if req.Type != models.MovementEntry && req.Type != models.MovementExit {
		return nil, errors.New("el tipo debe ser 'entry' o 'exit'")
	}

	product, err := s.productRepo.FindById(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("producto no encontrado: %w", err)
	}

	switch req.Type {
	case models.MovementEntry:
		product.Quantity += req.Quantity
	case models.MovementExit:
		if product.Quantity < req.Quantity {
			return nil, fmt.Errorf("stock insuficiente: disponible %d, solicitado %d", product.Quantity, req.Quantity)
		}
		product.Quantity -= req.Quantity
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("error al actualizar stock del producto: %w", err)
	}

	movement := &models.Movement{
		Product: models.ProductCategory{
			ID:   product.ID,
			Name: product.Name,
		},
		Type:     req.Type,
		Quantity: req.Quantity,
		Reason:   req.Reason,
		Date:     time.Now(),
	}

	if err := s.movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("error al registrar movimiento: %w", err)
	}

	return movement, nil
}

func (s *MovementService) FindAllMovements(ctx context.Context) ([]models.Movement, error) {
	movements, err := s.movementRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos: %w", err)
	}
	return movements, nil
}

func (s *MovementService) FindMovementsByProductID(ctx context.Context, productID string) ([]models.Movement, error) {
	_, err := s.productRepo.FindById(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("producto no encontrado: %w", err)
	}

	movements, err := s.movementRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos del producto: %w", err)
	}
	return movements, nil
}

func (s *MovementService) FindMovementsByType(ctx context.Context, movementType models.MovementType) ([]models.Movement, error) {
	if movementType != models.MovementEntry && movementType != models.MovementExit {
		return nil, errors.New("el tipo debe ser 'entry' o 'exit'")
	}

	movements, err := s.movementRepo.FindByType(ctx, movementType)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos por tipo: %w", err)
	}
	return movements, nil
}

func (s *MovementService) FindMovementsByDateRange(ctx context.Context, from, to time.Time) ([]models.Movement, error) {
	if from.After(to) {
		return nil, errors.New("la fecha inicial no puede ser posterior a la fecha final")
	}

	movements, err := s.movementRepo.FindByDateRange(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("error al obtener movimientos por rango de fecha: %w", err)
	}
	return movements, nil
}
