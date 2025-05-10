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
