CREATE DATABASE restaurants;
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