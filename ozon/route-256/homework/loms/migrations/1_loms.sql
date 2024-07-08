CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE order_items (
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sku_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, sku_id)
);

CREATE TABLE stocks (
    sku BIGINT PRIMARY KEY,
    total_count BIGINT NOT NULL,
    reserved BIGINT NOT NULL DEFAULT 0
);
