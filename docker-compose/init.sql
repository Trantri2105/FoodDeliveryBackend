CREATE DATABASE kong;
CREATE DATABASE users;
\c users
CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    email TEXT UNIQUE,
    password TEXT,
    name TEXT,
    gender TEXT,
    phone TEXT,
    role TEXT
);
INSERT INTO users (email, password, name, gender, phone, role) VALUES ('admin@gmail.com','$2a$10$Ac9T0McpRrLfF1KraoPgsOl/5r8qlbjcCtDanWyH.A4y7YyVi836G', 'admin', 'male', '1234567890','admin');

CREATE DATABASE restaurants;
\c restaurants

CREATE TABLE restaurants (
                             id SERIAL PRIMARY KEY,
                             name TEXT NOT NULL,
                             description TEXT,
                             address TEXT NOT NULL,
                             phone_number TEXT NOT NULL,
                             is_active BOOLEAN DEFAULT true,
                             open_time TIME NOT NULL,
                             close_time TIME NOT NULL
);

CREATE TABLE menu_items (
                            id SERIAL PRIMARY KEY,
                            restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
                            name TEXT NOT NULL,
                            description TEXT,
                            price INT NOT NULL,
                            is_available BOOLEAN DEFAULT true
);


-- Insert a restaurant
INSERT INTO restaurants (
    name,
    description,
    address,
    phone_number,
    is_active,
    open_time,
    close_time
) VALUES (
             'Coastal Breeze Bistro',
             'A charming seaside restaurant offering fresh seafood and Mediterranean-inspired dishes in a relaxed atmosphere with ocean views.',
             'Chung cư Season Avenue',
             '+1-415-555-8723',
             true,
             '11:00',
             '23:00'
         );

INSERT INTO menu_items (
    restaurant_id,
    name,
    description,
    price,
    is_available
) VALUES
-- Main dishes
(1, 'Pan-Seared Salmon', 'Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.', 24000, true),
(1, 'Seafood Linguine', 'Homemade linguine pasta tossed with shrimp, scallops, clams, and mussels in a light white wine garlic sauce.', 27000, true),
(1, 'Mediterranean Lamb Rack', 'Herb-crusted lamb rack served with rosemary potatoes and mint jelly.', 32000, true),
(1, 'Vegetable Paella', 'Saffron-infused rice with seasonal vegetables, artichoke hearts, and roasted peppers.', 19000, true),

-- Appetizers
(1, 'Calamari Fritti', 'Lightly battered calamari served with house-made marinara sauce and lemon wedges.', 15000, true),
(1, 'Bruschetta Sampler', 'Assortment of three bruschetta varieties: classic tomato and basil, olive tapenade, and goat cheese with honey.', 13000, false),
(1, 'Charcuterie Board', 'Selection of premium cured meats, artisanal cheeses, olives, nuts, and house-made bread.', 22000, true);

CREATE DATABASE orders;
\c orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT,
    shipping_address TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    status VARCHAR(30) NOT NULL,
    subtotal INT NOT NULL,
    delivery_fee INT NOT NULL,
    total_amount INT NOT NULL
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price INT NOT NULL,
    total_price INT NOT NULL
);

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
    status VARCHAR(32) CHECK (status IN ('pending', 'assigned', 'delivering', 'completed', 'canceled')) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);