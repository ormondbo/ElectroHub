-- ==========================================
-- ADDRESSES
-- ==========================================

CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(255) NOT NULL,

    CONSTRAINT chk_address_country
        CHECK (length(trim(country)) > 0),
    CONSTRAINT chk_address_city
        CHECK (length(trim(city)) > 0),
    CONSTRAINT chk_address_street
        CHECK (length(trim(street)) > 0)
);

-- ==========================================
-- CLIENTS
-- ==========================================

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_name VARCHAR(100) NOT NULL,
    client_surname VARCHAR(100) NOT NULL,
    birthday DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    address_id UUID NOT NULL,

    CONSTRAINT fk_client_address
        FOREIGN KEY (address_id)
        REFERENCES addresses(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_client_name
        CHECK (length(trim(client_name)) > 0),
    CONSTRAINT chk_client_surname
        CHECK (length(trim(client_surname)) > 0),
    CONSTRAINT chk_client_gender
        CHECK (gender IN ('male', 'female', 'other'))
);

-- ==========================================
-- SUPPLIERS
-- ==========================================

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address_id UUID NOT NULL,
    phone_number VARCHAR(20) NOT NULL,

    CONSTRAINT fk_supplier_address
        FOREIGN KEY (address_id)
        REFERENCES addresses(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_supplier_name
        CHECK (length(trim(name)) > 0),
    CONSTRAINT chk_supplier_phone
        CHECK (length(trim(phone_number)) > 0)
);

-- ==========================================
-- IMAGES
-- ==========================================

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image BYTEA NOT NULL,

    CONSTRAINT chk_image_not_empty
        CHECK (octet_length(image) > 0)
);

-- ==========================================
-- PRODUCTS
-- ==========================================

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price NUMERIC(12, 2) NOT NULL,
    available_stock INTEGER NOT NULL DEFAULT 0,
    last_update_date TIMESTAMPTZ,
    supplier_id UUID NOT NULL,
    image_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',

    CONSTRAINT fk_product_supplier
        FOREIGN KEY (supplier_id)
        REFERENCES suppliers(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_product_name
        CHECK (length(trim(name)) > 0),
    CONSTRAINT chk_product_category
        CHECK (length(trim(category)) > 0),
    CONSTRAINT chk_product_price
        CHECK (price >= 0),
    CONSTRAINT chk_product_stock
        CHECK (available_stock >= 0)
);

-- ==========================================
-- INDEXES
-- ==========================================

CREATE INDEX IF NOT EXISTS idx_clients_full_name
    ON clients(client_name, client_surname);

CREATE INDEX IF NOT EXISTS idx_products_supplier_id
    ON products(supplier_id);

CREATE INDEX IF NOT EXISTS idx_products_available
    ON products(id)
    WHERE available_stock > 0;

CREATE UNIQUE INDEX IF NOT EXISTS uq_product_image_not_empty
    ON products(image_id)
    WHERE image_id <> '00000000-0000-0000-0000-000000000000';

CREATE INDEX IF NOT EXISTS idx_clients_address_id
    ON clients(address_id);

CREATE INDEX IF NOT EXISTS idx_suppliers_address_id
    ON suppliers(address_id);