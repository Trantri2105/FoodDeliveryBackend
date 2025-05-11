CREATE DATABASE deliveries;
\c deliveries
CREATE TABLE shippers (
                          id INTEGER NOT NULL PRIMARY KEY,
                          password TEXT NOT NULL,
                          email TEXT UNIQUE NOT NULL,
                          role TEXT NOT NULL,
                          name TEXT NOT NULL,
                          gender TEXT CHECK (gender IN ('male', 'female', 'other')),
                          phone TEXT UNIQUE NOT NULL,
                          vehicle_type TEXT,
                          vehicle_plate TEXT,
                          total_deliveries INTEGER DEFAULT 0,
                          status TEXT CHECK (status IN ('available', 'unavailable', 'delivering', 'assigned')) NOT NULL
);
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
                            status VARCHAR(32) CHECK (status IN ('pending', 'assigned', 'delivering', 'delivered', 'canceled')) NOT NULL,
                            created_at TIMESTAMP DEFAULT NOW(),
                            updated_at TIMESTAMP DEFAULT NOW()
);