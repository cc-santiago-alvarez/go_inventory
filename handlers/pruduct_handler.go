package handlers

import (
	"net/http"

	"github.com/cc-santiago-alvarez/go_inventory.git/server"
	"github.com/cc-santiago-alvarez/go_inventory.git/services"
)

type ProductHandler struct {
	ProductService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (p *ProductHandler) CreateProductHandler(c *server.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		CategoryID  string  `json:"category_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		ResponseError(c, NewAppError("Datos invalidos", http.StatusBadRequest))
		return
	}

	if req.Name == "" || req.Price <= 0 {
		ResponseError(c, NewAppError("El nombre, el precio y la cantidad son requeridos", http.StatusBadRequest))
		return
	}

	product, err := p.ProductService.CreateProduct(c.Context(), req.Name, req.Description, req.Price, req.Quantity, req.CategoryID)

	if err != nil {
		ResponseError(c, NewAppError("Error al crear el producto", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Producto creado correctamente",
		"product": product,
	})
}

func (p *ProductHandler) GetAllProductsHandler(c *server.Context) {
	products, err := p.ProductService.FindAllProducts(c.Context())
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener los productos", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, products)
}

func (p *ProductHandler) GetProductByIdHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	product, err := p.ProductService.FindProductById(c.Context(), id)
	if err != nil {
		ResponseError(c, NewAppError("Error al obtener el producto", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, product)
}

func (p *ProductHandler) UpdateProductHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		CategoryID  string  `json:"category_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		ResponseError(c, NewAppError("Datos invalidos", http.StatusBadRequest))
		return
	}

	if req.Name == "" || req.Price <= 0 || req.Quantity <= 0 {
		ResponseError(c, NewAppError("El nombre, el precio y la cantidad son requeridos", http.StatusBadRequest))
		return
	}

	if err := p.ProductService.UpdateProduct(c.Context(), id, req.Name, req.Description, req.Price, req.Quantity, req.CategoryID); err != nil {
		ResponseError(c, NewAppError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Producto actualizado correctamente",
	})
}

func (p *ProductHandler) DeleteProductHandler(c *server.Context) {
	id := c.Request.PathValue("id")
	if id == "" {
		ResponseError(c, NewAppError("ID invalido", http.StatusBadRequest))
		return
	}

	if err := p.ProductService.DeleteProduct(c.Context(), id); err != nil {
		ResponseError(c, NewAppError(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Producto eliminado correctamente",
	})
}
