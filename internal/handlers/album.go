package handlers

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// GetAllAlbums responds with all albums in JSON format.
func GetAllAlbums(c echo.Context) error {
	albums, err := data.AllAlbums()
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, albums)
}

// GetAlbumsByArtist responds with albums filtered by artist name.
func GetAlbumsByArtist(c echo.Context) error {
	name := strings.TrimSpace(c.Param("name"))
	if name == "" {
		return c.JSON(400, map[string]string{"error": "Artist name is required"})
	}

	albums, err := data.AlbumsByArtist(name)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, albums)
}

// GetAlbumByID responds with a single album by its ID.
func GetAlbumByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(400, map[string]string{"error": "Invalid album ID"})
	}

	album, err := data.AlbumByID(id)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "Album not found"})
	}
	return c.JSON(200, album)
}

// CreateAlbum handles adding a new album to the database.
func CreateAlbum(c echo.Context) error {
	var album models.Album
	if err := c.Bind(&album); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid JSON"})
	}

	album.Title = strings.TrimSpace(album.Title)
	album.Artist = strings.TrimSpace(album.Artist)

	if album.Title == "" || album.Artist == "" {
		return c.JSON(400, map[string]string{"error": "Title and Artist are required"})
	}
	if len(album.Title) > 200 || len(album.Artist) > 100 {
		return c.JSON(400, map[string]string{"error": "Title or Artist too long"})
	}

	id, err := data.AddAlbum(album)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	album.ID = id
	return c.JSON(201, album)
}

// CanPurchaseAlbum checks if the requested quantity can be purchased.
func CanPurchaseAlbum(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(400, map[string]string{"error": "Invalid album ID"})
	}

	qtyStr := c.QueryParam("qty")
	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if err != nil || qty <= 0 {
		return c.JSON(400, map[string]string{"error": "Quantity must be positive"})
	}

	ok, err := data.CanPurchase(id, qty)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]bool{"canPurchase": ok})
}

// GetOrdersByUser serves the last 10 orders for a logged-in user (HTML page).
func GetOrdersByUser(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	auth, authOk := session.Values["authenticated"].(bool)
	userID, idOk := session.Values["user_id"].(int64)

	if !authOk || !auth || !idOk {
		session.Values["redirect_after_login"] = "/orders"
		session.Save(c.Request(), c.Response())
		c.Response().WriteHeader(401)
		tmpl := template.Must(template.ParseFiles("templates/login_required.html"))
		return tmpl.Execute(c.Response(), map[string]string{"Resource": "Orders API"})
	}

	cacheKey := fmt.Sprintf("orders:user:%d:last10", userID)
	var orders []models.GetOrder

	if err := data.GetOrdersCache(cacheKey, &orders); err != nil {
		log.Printf("Cache MISS for user: %d", userID)
		log.Printf("[SIMULATION] Sleeping 2s to simulate slow DB query for user %d...", userID)

		orders, err = data.GetOrdersByUser(userID)
		if err != nil {
			return c.JSON(500, map[string]string{"error": err.Error()})
		}

		_ = data.SetOrdersCache(cacheKey, orders)
	} else {
		log.Printf("Cache HIT for user: %d", userID)
	}

	c.Response().Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/orders.html"))
	return tmpl.Execute(c.Response(), map[string]any{"Orders": orders})
}

// CreateOrderByUser handles creating a new order for the logged-in user.
func CreateOrderByUser(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")
	userID, ok := session.Values["user_id"].(int64)
	auth, authOk := session.Values["authenticated"].(bool)

	if !ok || !authOk || !auth {
		session.Values["redirect_after_login"] = "/orders"
		session.Save(c.Request(), c.Response())
		return c.String(401, "Unauthorized: You must log in first to create an order.\nVisit /login to log in and retry.")
	}

	var order models.OrderRequest
	if err := c.Bind(&order); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid JSON"})
	}

	if order.AlbumID <= 0 || order.Quantity <= 0 {
		return c.JSON(400, map[string]string{"error": "AlbumID and Quantity must be positive"})
	}

	order.Customer = userID

	id, err := data.CreateOrderByUser(c.Request().Context(), order.AlbumID, order.Quantity, order.Customer)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	cacheKey := fmt.Sprintf("orders:user:%d:last10", order.Customer)
	var cached []models.GetOrder
	if err := data.GetOrdersCache(cacheKey, &cached); err == nil {
		newOrder := models.GetOrder{
			ID:       id,
			AlbumID:  order.AlbumID,
			Customer: order.Customer,
			Quantity: order.Quantity,
			Date:     time.Now(),
		}
		cached = append([]models.GetOrder{newOrder}, cached...)
		if len(cached) > 10 {
			cached = cached[:10]
		}
		_ = data.SetOrdersCache(cacheKey, cached)
		log.Printf("[CACHE UPDATE] User %d cache updated with new order %d", newOrder.Customer, newOrder.ID)
	} else {
		_ = data.OrderCache.Delete(cacheKey)
	}

	return c.JSON(201, map[string]any{
		"order_id": id,
		"message":  "Order created successfully",
	})
}

// GetCustomerName returns the full name of a customer by ID.
func GetCustomerName(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.JSON(400, map[string]string{"error": "Invalid customer ID"})
	}

	name, err := data.GetCustomerName(id)
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]string{"name": name})
}

// HandleMultipleResultSets demonstrates fetching multiple result sets (albums + customers).
func HandleMultipleResultSets(c echo.Context) error {
	result, err := data.GetAlbumsAndCustomers()
	if err != nil {
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, result)
}

// QueryWithTimeout executes a DB query with a timeout context.
func QueryWithTimeout(c echo.Context) error {
	albums, err := data.QueryAlbumsWithTimeout(c.Request().Context())
	if err != nil {
		return c.JSON(504, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, albums)
}
