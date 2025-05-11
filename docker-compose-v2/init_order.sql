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