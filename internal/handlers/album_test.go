package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"

	_ "modernc.org/sqlite"
)

// local test DB helper (self-contained for handlers tests)
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:") // in-memory DB
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE album (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT, artist TEXT, price REAL, quantity INTEGER
	)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO album (title, artist, price, quantity)
		VALUES ("Go Beats", "Gopher", 9.99, 5)`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setupAlbumHandlerDB(t *testing.T) {
	db := setupTestDB(t)
	data.InitDBConnection(db)
}

func TestGetAlbumByID(t *testing.T) {
	setupAlbumHandlerDB(t)

	// Create Echo instance
	e := echo.New()

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/albums/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set path parameter (Echo style)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call handler
	if err := GetAlbumByID(c); err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var alb models.Album
	if err := json.NewDecoder(rec.Body).Decode(&alb); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if strings.ToLower(alb.Artist) != "gopher" {
		t.Errorf("expected artist Gopher, got %s", alb.Artist)
	}
}
