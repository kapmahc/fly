CREATE TABLE mall_addresses (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  zip VARCHAR(12) NOT NULL,
  street VARCHAR(255) NOT NULL,
  city VARCHAR(32) NOT NULL,
  state VARCHAR(32) NOT NULL,
  country VARCHAR(32) NOT NULL,
  phone VARCHAR(32) NOT NULL,
  user_id BIGINT REFERENCES users,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_mall_addresses_usersname ON mall_addresses(username);
CREATE INDEX idx_mall_addresses_zip ON mall_addresses(zip);
CREATE INDEX idx_mall_addresses_city ON mall_addresses(city);
CREATE INDEX idx_mall_addresses_state ON mall_addresses(state);
CREATE INDEX idx_mall_addresses_country ON mall_addresses(country);
CREATE INDEX idx_mall_addresses_phone ON mall_addresses(phone);

CREATE TABLE mall_stores (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(8) NOT NULL,
  description TEXT NOT NULL,
  address_id BIGINT REFERENCES mall_addresses,
  owner_id BIGINT REFERENCES users,
  currency CHAR(3) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX idx_mall_stores_name ON mall_stores(name);
CREATE INDEX idx_mall_stores_currency ON mall_stores(currency);

CREATE TABLE mall_tags (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(8) NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE mall_vendors (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(8) NOT NULL,
  description TEXT NOT NULL,
  stores_id BIGINT REFERENCES mall_stores,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE mall_products (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(8) NOT NULL,
  description TEXT NOT NULL,
  vendor_id BIGINT REFERENCES mall_vendors,
  stores_id BIGINT REFERENCES mall_stores,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE mall_products_tags (
  tag_id BIGINT REFERENCES mall_tags ON DELETE CASCADE,
  product_id     BIGINT REFERENCES mall_products ON DELETE CASCADE,
  PRIMARY KEY (product_id, tag_id)
);

CREATE TABLE mall_variants(
  id BIGSERIAL PRIMARY KEY,
  sku VARCHAR(64) NOT NULL,
  product_id BIGINT REFERENCES mall_products,
  price NUMERIC(12,2) NOT NULL,
  cost NUMERIC(12,2) NOT NULL,
  weight NUMERIC(12,2) NOT NULL,
  height NUMERIC(12,2) NOT NULL,
  width NUMERIC(12,2) NOT NULL,
  length NUMERIC(12,2) NOT NULL,
  stores_id BIGINT REFERENCES mall_stores,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_mall_variants_sku ON mall_variants (sku);


CREATE TABLE mall_journals (
  id BIGSERIAL PRIMARY KEY,
  action VARCHAR(255) NOT NULL,
  quantity BIGINT NOT NULL,
  store_id BIGINT REFERENCES mall_stores,
  variant_id  BIGINT REFERENCES mall_variants,
  user_id BIGINT REFERENCES users,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE mall_stocks (
  id BIGSERIAL PRIMARY KEY,
  quantity BIGINT NOT NULL,
  store_id BIGINT REFERENCES mall_stores,
  variant_id  BIGINT REFERENCES mall_variants,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE mall_properties (
  id BIGSERIAL PRIMARY KEY,
  key VARCHAR(255) NOT NULL,
  val VARCHAR(2048) NOT NULL,
  variant_id BIGINT REFERENCES mall_variants,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_mall_properties_key_variant ON mall_properties (key, variant_id);
CREATE INDEX idx_mall_properties_key ON mall_properties (key);
