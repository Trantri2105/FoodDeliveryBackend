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
)