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
                            image_url TEXT,
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

INSERT INTO menu_items (restaurant_id, name, description, price, is_available, image_url) VALUES
                                                                                              (1, 'Pan-Seared Salmon', 'Fresh Atlantic salmon seared to perfection, served with roasted vegetables and lemon-dill sauce.', 24000, true, 'https://static01.nyt.com/images/2024/02/13/multimedia/LH-pan-seared-salmon-lwzt/LH-pan-seared-salmon-lwzt-mediumSquareAt3X.jpg'),
                                                                                              (1, 'Mediterranean Pasta', 'Al dente pasta with sun-dried tomatoes, olives, feta cheese, and fresh herbs in a light olive oil sauce.', 27000, true, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSrNGbsLrQaN7TyYv93p2gGhyDKcd9MzqshkA&s'),
                                                                                              (1, 'Grilled Seafood Platter', 'An assortment of grilled shrimp, scallops, and fish fillets, served with a citrus butter sauce and seasonal vegetables.', 32000, true, 'https://eastsidebarandgrill.com.au/wp-content/uploads/2022/03/Eastside_Seafood-Platter_1200x800px.jpg'),
                                                                                              (1, 'Lemon Herb Chicken', 'Tender chicken breast marinated in lemon and herbs, grilled to perfection and served with roasted potatoes.', 22000, true, 'https://www.foodandwine.com/thmb/t9YqzGbmH-huAbV6xitCQs0-G4s=/1500x0/filters:no_upscale():max_bytes(150000):strip_icc()/FAW-recipes-herb-and-lemon-roasted-chicken-hero-c4ba0aec56884683be482c47b1e1df11.jpg'),
                                                                                              (1, 'Vegetarian Risotto', 'Creamy arborio rice cooked with wild mushrooms, asparagus, and parmesan cheese.', 20000, true, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcR6dgoSlQmnXyuVWdQUZU3ERZAgbXmXGtlU6g&s');