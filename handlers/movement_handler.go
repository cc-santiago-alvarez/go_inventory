package handlers

import (
	"net/http"

	"github.com/cc-santiago-alvarez/go_inventory.git/models"
	"github.com/cc-santiago-alvarez/go_inventory.git/server"
	"github.com/cc-santiago-alvarez/go_inventory.git/services"
)

type MovementHandler struct {
	MovementService *services.MovementService
}

func NewMovementHandler(movementService *services.MovementService) *MovementHandler {
	return &MovementHandler{MovementService: movementService}
}

func (m *MovementHandler) CreateMovementHandler(c *server.Context) {
	var req models.CreateMovementRequest

	if err := c.BindJSON(&req); err != nil {
		ResponseError(c, NewAppError("Datos invalidos", http.StatusBadRequest))
		return
	}

	if req.ProductID == "" || req.Quantity <= 0 || req.Reason == "" {
		ResponseError(c, NewAppError("product_id, quantity y reason son requeridos", http.StatusBadRequest))
		return
	}

	if req.Type != models.MovementEntry && req.Type != models.MovementExit {
		ResponseError(c, NewAppError("El tipo debe ser 'entry' o 'exit'", http.StatusBadRequest))
		return
	}

	movement, err := m.MovementService.CreateMovement(c.Context(), req)
	if err != nil {
		ResponseError(c, NewAppError(err.Error(), http.StatusBadRequest))
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Movimiento registrado correctamente",
		"movement": movement,
	})
}

func (m *MovementHandler) GetAllMovementsHandler(c *server.Context) {
	movements, err := m.MovementService.FindAllMovements(c.Context())
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener los movimientos", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, movements)
}

func (m *MovementHandler) GetMovementsByProductHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	movements, err := m.MovementService.FindMovementsByProductID(c.Context(), id)
	if err != nil {
		ResponseError(c, NewAppError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, movements)
}
