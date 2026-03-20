package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const schema = `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    prefix VARCHAR(10) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    price DECIMAL(10, 2) NOT NULL,
    quantity INT,
    category_id UUID NOT NULL REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	_ = godotenv.Load(".env")

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL o DATABASE_URL debe estar configurada para ejecutar los tests de integración")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos de test: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Error al hacer ping a la base de datos de test: %v", err)
	}

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("Error al ejecutar el schema: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("TRUNCATE products, categories CASCADE")
    db.Close()
	})

	fmt.Println("Base de datos de test configurada correctamente")
	return db
}
