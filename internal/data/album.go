package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	rows, err := db.Query("SELECT id, title, artist, price, quantity FROM albums")
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
	rows, err := db.Query("SELECT id, title, artist, price, quantity FROM albums WHERE artist = ?", name)
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
	err := db.QueryRow("SELECT id, title, artist, price, quantity FROM albums WHERE id = ?", id).
		Scan(&album.ID, &album.Title, &album.Artist, &album.Price, &album.Quantity)
	if err != nil {
		return album, err
	}
	return album, nil
}

// AddAlbum inserts a new album and returns its inserted ID.
func AddAlbum(alb models.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO albums (title, artist, price, quantity) VALUES (?, ?, ?, ?)",
		alb.Title, alb.Artist, alb.Price, alb.Quantity)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// CanPurchase checks if the requested quantity is available for a given album.
func CanPurchase(id int64, quantity int64) (bool, error) {
	var enough bool
	err := db.QueryRow("SELECT (quantity >= ?) FROM albums WHERE id = ?", quantity, id).Scan(&enough)
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
		FROM album_orders
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
	if err := tx.QueryRowContext(ctx, "SELECT (quantity >= ?) FROM albums WHERE id = ?", quantity, albumID).Scan(&enough); err != nil {
		return 0, err
	}
	if !enough {
		return 0, fmt.Errorf("not enough inventory")
	}

	if _, err := tx.ExecContext(ctx, "UPDATE albums SET quantity = quantity - ? WHERE id = ?", quantity, albumID); err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, "INSERT INTO album_orders (album_id, cust_id, quantity, date) VALUES (?, ?, ?, ?)",
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

// GetAlbumsAndCustomers returns albums and customers in a combined map using multiple result sets.
func GetAlbumsAndOrders() (map[string]any, error) {
	rows, err := db.Query("SELECT * FROM albums; SELECT * FROM album_orders;")
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

	var orders []models.GetOrder
	if rows.NextResultSet() {
		var o models.GetOrder
		for rows.Next() {
			if err := rows.Scan(&o.ID, &o.AlbumID, &o.Customer, &o.Quantity, &o.Date); err != nil {
				return nil, err
			}
			orders = append(orders, o)
		}
	}

	return map[string]any{
		"albums":       albums,
		"album_orders": orders,
	}, nil
}

// QueryAlbumsWithTimeout queries albums with a context timeout.
func QueryAlbumsWithTimeout(ctx context.Context) ([]models.Album, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Println("⏳ Starting album query with a 5-second timeout...")
	log.Println("🐢 Simulating slow query with SLEEP(10) to trigger timeout...")

	rows, err := db.QueryContext(ctx, "SELECT SLEEP(10), id, title, artist, price, quantity FROM albums")

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("⚠️  Album query timed out")
		} else {
			log.Printf("❌ Album query failed: %v\n", err)
		}
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

	log.Printf("✅ Album query finished, fetched %d rows\n", len(albums))
	return albums, rows.Err()
}
