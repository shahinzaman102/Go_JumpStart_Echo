CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    price DECIMAL(5,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0
);

INSERT INTO albums (id, title, artist, price, quantity)
VALUES
    (1, 'Blue Train', 'John Coltrane', 56.99, 10),
    (2, 'Giant Steps', 'John Coltrane', 63.99, 8),
    (3, 'Jeru', 'Gerry Mulligan', 17.99, 12),
    (4, 'Sarah Vaughan', 'Sarah Vaughan', 34.98, 5)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS album_orders (
    id SERIAL PRIMARY KEY,
    album_id INT NOT NULL,
    cust_id INT NOT NULL,
    quantity INT NOT NULL,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (album_id) REFERENCES albums(id)
);

-- Synchronize the albums auto-increment sequence with the
-- current maximum ID to prevent duplicate primary key errors
-- when new records are inserted after manually assigning IDs.
SELECT setval(
    pg_get_serial_sequence('albums', 'id'),
    COALESCE((SELECT MAX(id) FROM albums), 1),
    EXISTS (SELECT 1 FROM albums)
);
