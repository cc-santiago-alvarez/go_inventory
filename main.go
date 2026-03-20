package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/cc-santiago-alvarez/go_inventory.git/graphql/explorer"
	"github.com/cc-santiago-alvarez/go_inventory.git/config"
	"github.com/cc-santiago-alvarez/go_inventory.git/database"
	"github.com/cc-santiago-alvarez/go_inventory.git/graphql/generated"
	"github.com/cc-santiago-alvarez/go_inventory.git/graphql/resolvers"
	"github.com/cc-santiago-alvarez/go_inventory.git/handlers"
	"github.com/cc-santiago-alvarez/go_inventory.git/repositories"
	"github.com/cc-santiago-alvarez/go_inventory.git/server"
	"github.com/cc-santiago-alvarez/go_inventory.git/services"
)

// func main() {
// 	config := config.LoadConfig()

// 	if err := database.Connect(config.DatabaseURL); err != nil {
// 		log.Fatal("Error al conectar con la base de datos", err)
// 	}

// 	defer database.Close()

// 	//Inicializar repositorios
// 	productRepo := repositories.NewProductRepository(database.DB)
// 	categoryRepo := repositories.NewCategoryRepository(database.DB)
// 	movementRepo := repositories.NewMovementRepository(database.DB)

// 	// Inicializar servicios
// 	productService := services.NewProductService(productRepo, categoryRepo)
// 	categoryService := services.NewCategoryService(categoryRepo, productRepo)
// 	movementService := services.NewMovementService(movementRepo, productRepo)

// 	//Inicializar handlers
// 	productHandler := handlers.NewProductHandler(productService)
// 	categoryHandler := handlers.NewCategoryHandler(categoryService)
// 	movementHandler := handlers.NewMovementHandler(movementService)

// 	app := server.NewApp()

// 	//rutas
// 	app.Get("/products", productHandler.GetAllProductsHandler)
// 	app.Get("/products/{id}", productHandler.GetProductByIdHandler)
// 	app.Post("/products", productHandler.CreateProductHandler)
// 	app.Put("/products/{id}", productHandler.UpdateProductHandler)
// 	app.Delete("/products/{id}", productHandler.DeleteProductHandler)

// 	// Rutas de categorías
// 	app.Get("/category", categoryHandler.GetAllCategoriesHandler)
// 	app.Get("/category/{id}", categoryHandler.GetCategoryByIdHandler)
// 	app.Post("/category", categoryHandler.CreateCategoryHandler)
// 	app.Put("/category/{id}", categoryHandler.UpdateCategoryHandler)
// 	app.Delete("/category/{id}", categoryHandler.DeleteCategoryHandler)

// 	app.Get("/category/{id}/products", categoryHandler.GetCategoryWithProductsHandler)

// 	// Rutas de movimientos
// 	app.Post("/movements", movementHandler.CreateMovementHandler)
// 	app.Get("/movements", movementHandler.GetAllMovementsHandler)
// 	app.Get("/movements/product/{id}", movementHandler.GetMovementsByProductHandler)

// 	if err := app.RunServer(config.Port); err != nil {
// 		log.Fatal("Error al iniciar servidor", err)
// 	}
// }

func main() {
	cfg := config.LoadConfig()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Error al conectar con la base de datos", err)
	}
	defer database.Close()

	productRepo := repositories.NewProductRepository(database.DB)
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	movementRepo := repositories.NewMovementRepository(database.DB)

	productService := services.NewProductService(productRepo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	movementService := services.NewMovementService(movementRepo, productRepo)

	app := server.NewApp()

	if cfg.API_MODE == "http" || cfg.API_MODE == "both" {
		registerHTTPRoutes(app, productService, categoryService, movementService)
	}

	if cfg.API_MODE == "graphql" || cfg.API_MODE == "both" {
		registerGraphQL(app, productService, categoryService, movementService)
	}

	if err := app.RunServer(cfg.Port); err != nil {
		log.Fatal("Error al iniciar servidor", err)
	}
}

func registerHTTPRoutes(app *server.App, ps *services.ProductService, cs *services.CategoryService, ms *services.MovementService) {
	productHandler := handlers.NewProductHandler(ps)
	categoryHandler := handlers.NewCategoryHandler(cs)
	movementHandler := handlers.NewMovementHandler(ms)

	app.Get("/products", productHandler.GetAllProductsHandler)
	app.Get("/products/{id}", productHandler.GetProductByIdHandler)
	app.Post("/products", productHandler.CreateProductHandler)
	app.Put("/products/{id}", productHandler.UpdateProductHandler)
	app.Delete("/products/{id}", productHandler.DeleteProductHandler)

	app.Get("/category", categoryHandler.GetAllCategoriesHandler)
	app.Get("/category/{id}", categoryHandler.GetCategoryByIdHandler)
	app.Post("/category", categoryHandler.CreateCategoryHandler)
	app.Put("/category/{id}", categoryHandler.UpdateCategoryHandler)
	app.Delete("/category/{id}", categoryHandler.DeleteCategoryHandler)
	app.Get("/category/{id}/products", categoryHandler.GetCategoryWithProductsHandler)

	app.Post("/movements", movementHandler.CreateMovementHandler)
	app.Get("/movements", movementHandler.GetAllMovementsHandler)
	app.Get("/movements/product/{id}", movementHandler.GetMovementsByProductHandler)
}

func registerGraphQL(app *server.App, ps *services.ProductService, cs *services.CategoryService, ms *services.MovementService) {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: &resolvers.Resolver{
				ProductService:  ps,
				CategoryService: cs,
				MovementService: ms,
			},
		}),
	)

	app.Handle("/graphql", srv)
	app.Handle("/playground", explorer.Handler("GraphQL Explorer", "/graphql"))
}
