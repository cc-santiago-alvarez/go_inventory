package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import "github.com/cc-santiago-alvarez/go_inventory.git/services"

type Resolver struct {
	ProductService  *services.ProductService
	CategoryService *services.CategoryService
	MovementService *services.MovementService
} 