CREATE TABLE deliveries (
    delivery_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    shipper_id INT NOT NULL REFERENCES shippers(id),
    restaurant_address TEXT NOT NULL,
    shipping_address TEXT NOT NULL,
    distance DOUBLE PRECISION NOT NULL,
    duration DOUBLE PRECISION NOT NULL,
    fee INT NOT NULL,
    from_coords JSONB NOT NULL,
    to_coords JSONB NOT NULL,
    geometry_line TEXT NOT NULL,
    status VARCHAR(32) CHECK (status IN ('pending', 'assigned', 'delivering', 'completed', 'canceled')) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
