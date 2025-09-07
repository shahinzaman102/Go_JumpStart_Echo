package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

var db *sql.DB

// InitDBConnection sets the package-level DB variable for reuse across data access functions.
func InitDBConnection(conn *sql.DB) {
	db = conn
}

// AllAlbums returns all albums in the database.
func AllAlbums() ([]models.Album, error) {
	rows, err := db.Query("SELECT id, title, artist, price, quantity FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price, &a.Quantity); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}

// AlbumsByArtist returns albums filtered by the artist's name.
func AlbumsByArtist(name string) ([]models.Album, error) {
	rows, err := db.Query("SELECT id, title, artist, price, quantity FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var alb models.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price, &alb.Quantity); err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}
	return albums, rows.Err()
}

// AlbumByID retrieves a single album by its ID.
func AlbumByID(id int64) (models.Album, error) {
	var album models.Album
	err := db.QueryRow("SELECT id, title, artist, price, quantity FROM album WHERE id = ?", id).
		Scan(&album.ID, &album.Title, &album.Artist, &album.Price, &album.Quantity)
	if err != nil {
		return album, err
	}
	return album, nil
}

// AddAlbum inserts a new album and returns its inserted ID.
func AddAlbum(alb models.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price, quantity) VALUES (?, ?, ?, ?)",
		alb.Title, alb.Artist, alb.Price, alb.Quantity)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// CanPurchase checks if the requested quantity is available for a given album.
func CanPurchase(id int64, quantity int64) (bool, error) {
	var enough bool
	err := db.QueryRow("SELECT (quantity >= ?) FROM album WHERE id = ?", quantity, id).Scan(&enough)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("unknown album ID %d", id)
		}
		return false, err
	}
	return enough, nil
}

// GetOrdersByUser returns the last 10 orders for a customer.
func GetOrdersByUser(userID int64) ([]models.GetOrder, error) {
	time.Sleep(2 * time.Second) // Artificial delay for testing only

	rows, err := db.Query(`
		SELECT id, album_id, cust_id, quantity, date
		FROM album_order
		WHERE cust_id = ?
		ORDER BY date DESC
		LIMIT 10
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.GetOrder
	for rows.Next() {
		var o models.GetOrder
		if err := rows.Scan(&o.ID, &o.AlbumID, &o.Customer, &o.Quantity, &o.Date); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

// CreateOrderByUser creates an order for a user within a transaction (all-or-nothing).
func CreateOrderByUser(ctx context.Context, albumID, quantity, custID int64) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var enough bool
	if err := tx.QueryRowContext(ctx, "SELECT (quantity >= ?) FROM album WHERE id = ?", quantity, albumID).Scan(&enough); err != nil {
		return 0, err
	}
	if !enough {
		return 0, fmt.Errorf("not enough inventory")
	}

	if _, err := tx.ExecContext(ctx, "UPDATE album SET quantity = quantity - ? WHERE id = ?", quantity, albumID); err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, "INSERT INTO album_order (album_id, cust_id, quantity, date) VALUES (?, ?, ?, ?)",
		albumID, custID, quantity, time.Now())
	if err != nil {
		return 0, err
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return orderID, nil
}

// GetCustomerName retrieves a customer's full name by ID.
func GetCustomerName(id int64) (string, error) {
	var name string
	if err := db.QueryRow("SELECT full_name FROM customer WHERE id = ?", id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("customer not found")
		}
		return "", err
	}
	return name, nil
}

// GetAlbumsAndCustomers returns albums and customers in a combined map using multiple result sets.
func GetAlbumsAndCustomers() (map[string]any, error) {
	rows, err := db.Query("SELECT * FROM album; SELECT * FROM customer;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price, &a.Quantity); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}

	var customers []map[string]any
	if rows.NextResultSet() {
		for rows.Next() {
			var id int64
			var fullName, address, phone string
			if err := rows.Scan(&id, &fullName, &address, &phone); err != nil {
				return nil, err
			}
			customers = append(customers, map[string]any{
				"id":       id,
				"fullName": fullName,
				"address":  address,
				"phone":    phone,
			})
		}
	}

	return map[string]any{
		"albums":    albums,
		"customers": customers,
	}, nil
}

// QueryAlbumsWithTimeout queries albums with a context timeout.
func QueryAlbumsWithTimeout(ctx context.Context) ([]models.Album, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT id, title, artist, price, quantity FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price, &a.Quantity); err != nil {
			return nil, err
		}
		albums = append(albums, a)
	}
	return albums, rows.Err()
}
