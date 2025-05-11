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