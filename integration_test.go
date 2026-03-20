package main

import (
	"context"
	"testing"

	"github.com/cc-santiago-alvarez/go_inventory.git/repositories"
	"github.com/cc-santiago-alvarez/go_inventory.git/services"
	"github.com/cc-santiago-alvarez/go_inventory.git/testutil"
)

func TestCreateCategory(t *testing.T) {
	db := testutil.SetupTestDB(t)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	ctx := context.Background()

	t.Run("crear categoría exitosamente", func(t *testing.T) {
		category, err := categoryService.CreateCategory(ctx, "Electrónica", "Productos electrónicos")
		if err != nil {
			t.Fatalf("No se esperaba error al crear categoría: %v", err)
		}

		if category.ID == "" {
			t.Error("Se esperaba un ID generado para la categoría")
		}
		if category.Name != "Electrónica" {
			t.Errorf("Se esperaba nombre 'Electrónica', se obtuvo '%s'", category.Name)
		}
		if category.Description != "Productos electrónicos" {
			t.Errorf("Se esperaba descripción 'Productos electrónicos', se obtuvo '%s'", category.Description)
		}

		t.Logf("Categoría creada con ID: %s", category.ID)
	})

	t.Run("error al crear categoría sin nombre", func(t *testing.T) {
		_, err := categoryService.CreateCategory(ctx, "", "Sin nombre")
		if err == nil {
			t.Error("Se esperaba error al crear categoría sin nombre")
		}
	})

	t.Run("error al crear categoría con nombre duplicado", func(t *testing.T) {
		_, err := categoryService.CreateCategory(ctx, "Electrónica", "Duplicada")
		if err == nil {
			t.Error("Se esperaba error al crear categoría con nombre duplicado")
		}
	})
}

func TestCreateProduct(t *testing.T) {
	db := testutil.SetupTestDB(t)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	ctx := context.Background()

	category, err := categoryService.CreateCategory(ctx, "Ropa", "Prendas de vestir")
	if err != nil {
		t.Fatalf("Error al crear categoría previa: %v", err)
	}

	t.Run("crear producto exitosamente", func(t *testing.T) {
		product, err := productService.CreateProduct(ctx, "Camiseta", "Camiseta de algodón", 29.99, 100, category.ID)
		if err != nil {
			t.Fatalf("No se esperaba error al crear producto: %v", err)
		}

		if product.ID == "" {
			t.Error("Se esperaba un ID generado para el producto")
		}
		if product.Name != "Camiseta" {
			t.Errorf("Se esperaba nombre 'Camiseta', se obtuvo '%s'", product.Name)
		}
		if product.Price != 29.99 {
			t.Errorf("Se esperaba precio 29.99, se obtuvo %f", product.Price)
		}
		if product.CategoryID != category.ID {
			t.Errorf("Se esperaba category_id '%s', se obtuvo '%s'", category.ID, product.CategoryID)
		}

		t.Logf("Producto creado con ID: %s, asociado a categoría: %s", product.ID, product.CategoryID)
	})

	t.Run("error al crear producto sin nombre", func(t *testing.T) {
		_, err := productService.CreateProduct(ctx, "", "Sin nombre", 10.0, 5, category.ID)
		if err == nil {
			t.Error("Se esperaba error al crear producto sin nombre")
		}
	})

	t.Run("error al crear producto con precio cero", func(t *testing.T) {
		_, err := productService.CreateProduct(ctx, "Producto", "Desc", 0, 5, category.ID)
		if err == nil {
			t.Error("Se esperaba error al crear producto con precio 0")
		}
	})

	t.Run("error al crear producto con cantidad cero", func(t *testing.T) {
		_, err := productService.CreateProduct(ctx, "Producto", "Desc", 10.0, 0, category.ID)
		if err == nil {
			t.Error("Se esperaba error al crear producto con cantidad 0")
		}
	})

	t.Run("error al crear producto sin category_id", func(t *testing.T) {
		_, err := productService.CreateProduct(ctx, "Producto", "Desc", 10.0, 5, "")
		if err == nil {
			t.Error("Se esperaba error al crear producto sin category_id")
		}
	})

	t.Run("error al crear producto con categoría inexistente", func(t *testing.T) {
		_, err := productService.CreateProduct(ctx, "Producto", "Desc", 10.0, 5, "00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Error("Se esperaba error al crear producto con categoría inexistente")
		}
	})
}

func TestCategoryProductRelationship(t *testing.T) {
	db := testutil.SetupTestDB(t)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	categoryService := services.NewCategoryService(categoryRepo, productRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	ctx := context.Background()

	catElectronica, err := categoryService.CreateCategory(ctx, "Electrónica", "Dispositivos electrónicos")
	if err != nil {
		t.Fatalf("Error al crear categoría Electrónica: %v", err)
	}

	catAlimentos, err := categoryService.CreateCategory(ctx, "Alimentos", "Productos alimenticios")
	if err != nil {
		t.Fatalf("Error al crear categoría Alimentos: %v", err)
	}

	laptop, err := productService.CreateProduct(ctx, "Laptop", "Laptop gaming", 1500.00, 10, catElectronica.ID)
	if err != nil {
		t.Fatalf("Error al crear producto Laptop: %v", err)
	}

	mouse, err := productService.CreateProduct(ctx, "Mouse", "Mouse inalámbrico", 25.00, 50, catElectronica.ID)
	if err != nil {
		t.Fatalf("Error al crear producto Mouse: %v", err)
	}

	arroz, err := productService.CreateProduct(ctx, "Arroz", "Arroz integral 1kg", 3.50, 200, catAlimentos.ID)
	if err != nil {
		t.Fatalf("Error al crear producto Arroz: %v", err)
	}

	t.Run("categoría Electrónica tiene 2 productos", func(t *testing.T) {
		catWithProducts, err := categoryService.FindCategoryWithProducts(ctx, catElectronica.ID)
		if err != nil {
			t.Fatalf("Error al obtener categoría con productos: %v", err)
		}

		if catWithProducts.Name != "Electrónica" {
			t.Errorf("Se esperaba nombre 'Electrónica', se obtuvo '%s'", catWithProducts.Name)
		}

		if len(catWithProducts.Products) != 2 {
			t.Fatalf("Se esperaban 2 productos en Electrónica, se obtuvieron %d", len(catWithProducts.Products))
		}

		productNames := map[string]bool{}
		for _, p := range catWithProducts.Products {
			productNames[p.Name] = true
			if p.CategoryID != catElectronica.ID {
				t.Errorf("Producto '%s' tiene category_id '%s', se esperaba '%s'", p.Name, p.CategoryID, catElectronica.ID)
			}
		}

		if !productNames["Laptop"] {
			t.Error("Se esperaba encontrar 'Laptop' en los productos de Electrónica")
		}
		if !productNames["Mouse"] {
			t.Error("Se esperaba encontrar 'Mouse' en los productos de Electrónica")
		}

		t.Logf("Electrónica contiene: Laptop y Mouse (%d productos)", len(catWithProducts.Products))
	})

	t.Run("categoría Alimentos tiene 1 producto", func(t *testing.T) {
		catWithProducts, err := categoryService.FindCategoryWithProducts(ctx, catAlimentos.ID)
		if err != nil {
			t.Fatalf("Error al obtener categoría con productos: %v", err)
		}

		if len(catWithProducts.Products) != 1 {
			t.Fatalf("Se esperaba 1 producto en Alimentos, se obtuvieron %d", len(catWithProducts.Products))
		}

		if catWithProducts.Products[0].Name != "Arroz" {
			t.Errorf("Se esperaba 'Arroz', se obtuvo '%s'", catWithProducts.Products[0].Name)
		}

		t.Logf("Alimentos contiene: Arroz (%d producto)", len(catWithProducts.Products))
	})

	t.Run("los productos tienen el category_id correcto", func(t *testing.T) {
		foundLaptop, err := productService.FindProductById(ctx, laptop.ID)
		if err != nil {
			t.Fatalf("Error al buscar Laptop: %v", err)
		}
		if foundLaptop.CategoryID != catElectronica.ID {
			t.Errorf("Laptop debería pertenecer a Electrónica (%s), pero tiene category_id '%s'", catElectronica.ID, foundLaptop.CategoryID)
		}

		foundMouse, err := productService.FindProductById(ctx, mouse.ID)
		if err != nil {
			t.Fatalf("Error al buscar Mouse: %v", err)
		}
		if foundMouse.CategoryID != catElectronica.ID {
			t.Errorf("Mouse debería pertenecer a Electrónica (%s), pero tiene category_id '%s'", catElectronica.ID, foundMouse.CategoryID)
		}

		foundArroz, err := productService.FindProductById(ctx, arroz.ID)
		if err != nil {
			t.Fatalf("Error al buscar Arroz: %v", err)
		}
		if foundArroz.CategoryID != catAlimentos.ID {
			t.Errorf("Arroz debería pertenecer a Alimentos (%s), pero tiene category_id '%s'", catAlimentos.ID, foundArroz.CategoryID)
		}

		t.Log("Todos los productos tienen el category_id correcto")
	})

	t.Run("categoría sin productos retorna lista vacía", func(t *testing.T) {
		catVacia, err := categoryService.CreateCategory(ctx, "Vacía", "Categoría sin productos")
		if err != nil {
			t.Fatalf("Error al crear categoría vacía: %v", err)
		}

		catWithProducts, err := categoryService.FindCategoryWithProducts(ctx, catVacia.ID)
		if err != nil {
			t.Fatalf("Error al obtener categoría vacía con productos: %v", err)
		}

		if len(catWithProducts.Products) != 0 {
			t.Errorf("Se esperaban 0 productos, se obtuvieron %d", len(catWithProducts.Products))
		}

		t.Log("Categoría vacía retorna correctamente lista de productos vacía")
	})

	t.Run("no se puede eliminar categoría con productos asociados", func(t *testing.T) {
		err := categoryService.DeleteCategory(ctx, catElectronica.ID)
		if err == nil {
			t.Error("Se esperaba error al eliminar categoría con productos asociados (restricción FK)")
		}
		t.Logf("Correctamente rechazó eliminar categoría con productos: %v", err)
	})
}
