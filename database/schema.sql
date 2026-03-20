CREATE EXTENSION IF NOT EXISTS "pgcrypto";
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    prefix VARCHAR(5) UNIQUE, --NOT NULL
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(15)  UNIQUE, --NOT NULL
    name VARCHAR(255) NOT NULL,
    description  VARCHAR(255),
    price DECIMAL(10, 2) NOT NULL,
    quantity INT,
    category_id UUID REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('entry', 'exit')),
    quantity INT NOT NULL CHECK (quantity > 0),
    reason VARCHAR(255), --NOT NULL
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);