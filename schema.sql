CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATETIME
);

CREATE TABLE IF NOT EXISTS album (
    id       INT AUTO_INCREMENT NOT NULL,
    title    VARCHAR(128) NOT NULL,
    artist   VARCHAR(255) NOT NULL,
    price    DECIMAL(5,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);

INSERT IGNORE INTO album (id, title, artist, price, quantity)
VALUES
    (1, 'Blue Train', 'John Coltrane', 56.99, 10),
    (2, 'Giant Steps', 'John Coltrane', 63.99, 8),
    (3, 'Jeru', 'Gerry Mulligan', 17.99, 12),
    (4, 'Sarah Vaughan', 'Sarah Vaughan', 34.98, 5);

CREATE TABLE IF NOT EXISTS customer (
    id INT AUTO_INCREMENT PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL
);

INSERT IGNORE INTO customer (id, full_name, address, phone)
VALUES
    (1, 'John Doe', '12/A, Dhanmondi, Dhaka, Bangladesh', '+88017xxxxxxx'),
    (2, 'Jane Smith', '45/B, Banani, Dhaka, Bangladesh', '+88019xxxxxxx'),
    (3, 'Michael Johnson', '78/C, Gulshan, Dhaka, Bangladesh', '+88018xxxxxxx'),
    (4, 'Emily Davis', '32/D, Chittagong, Bangladesh', '+88016xxxxxxx'),
    (5, 'David Wilson', '65/E, Khulna, Bangladesh', '+88015xxxxxxx');

CREATE TABLE IF NOT EXISTS album_order (
    id INT AUTO_INCREMENT PRIMARY KEY,
    album_id INT NOT NULL,
    cust_id INT NOT NULL,
    quantity INT NOT NULL,
    date DATETIME NOT NULL,
    FOREIGN KEY (album_id) REFERENCES album(id),
    FOREIGN KEY (cust_id) REFERENCES customer(id)
);
