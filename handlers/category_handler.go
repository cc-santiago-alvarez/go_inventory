package handlers

import (
	"net/http"

	"github.com/cc-santiago-alvarez/go_inventory.git/server"
	"github.com/cc-santiago-alvarez/go_inventory.git/services"
)

type CategoryHandler struct {
	CategoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: categoryService}
}

func (p *CategoryHandler) CreateCategoryHandler(c *server.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		ResponseError(c, NewAppError("Datos invalidos", http.StatusBadRequest))
		return
	}

	if req.Name == "" {
		ResponseError(c, NewAppError("El nombre de la categoria es requerido", http.StatusBadRequest))
		return
	}

	category, err := p.CategoryService.CreateCategory(c.Context(), req.Name, req.Description)
	if err != nil {
		ResponseError(c, NewAppError("Error al crear la categoria", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Categoria creada correctamente",
		"category": category,
	})
}

func (p *CategoryHandler) GetAllCategoriesHandler(c *server.Context) {
	categories, err := p.CategoryService.FindAllCategories(c.Context())
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener las categorias", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (p *CategoryHandler) GetCategoryByIdHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	category, err := p.CategoryService.FindCategoryById(c.Context(), id)
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener la categoria", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, category)
}

func (p *CategoryHandler) UpdateCategoryHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&req); err != nil {
		ResponseError(c, NewAppError("Datos invalidos", http.StatusBadRequest))
		return
	}

	if req.Name == "" {
		ResponseError(c, NewAppError("El nombre de la categoria es requerido", http.StatusBadRequest))
		return
	}

	if err := p.CategoryService.UpdateCategory(c.Context(), id, req.Name, req.Description); err != nil {
		ResponseError(c, NewAppError("Error al actualizar la categoria", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Categoria actualizada correctamente",
	})
}

func (p *CategoryHandler) DeleteCategoryHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	if err := p.CategoryService.DeleteCategory(c.Context(), id); err != nil {
		ResponseError(c, NewAppError("Error al eliminar la categoria", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Categoria eliminada correctamente",
	})
}

func (p *CategoryHandler) GetCategoryWithProductsHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	category, err := p.CategoryService.FindCategoryWithProducts(c.Context(), id)
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener la categoría con productos", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, category)
}
