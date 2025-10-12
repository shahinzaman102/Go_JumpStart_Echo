CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME
);

CREATE TABLE IF NOT EXISTS albums (
    id       INT AUTO_INCREMENT NOT NULL,
    title    VARCHAR(128) NOT NULL,
    artist   VARCHAR(255) NOT NULL,
    price    DECIMAL(5,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);

INSERT IGNORE INTO albums (id, title, artist, price, quantity)
VALUES
    (1, 'Blue Train', 'John Coltrane', 56.99, 10),
    (2, 'Giant Steps', 'John Coltrane', 63.99, 8),
    (3, 'Jeru', 'Gerry Mulligan', 17.99, 12),
    (4, 'Sarah Vaughan', 'Sarah Vaughan', 34.98, 5);

CREATE TABLE IF NOT EXISTS album_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    album_id INT NOT NULL,
    cust_id INT NOT NULL,
    quantity INT NOT NULL,
    date DATETIME NOT NULL,
    FOREIGN KEY (album_id) REFERENCES albums(id)
);
