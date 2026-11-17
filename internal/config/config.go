package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Store *sessions.CookieStore

// InitEnv loads .env file so os.Getenv() works.
func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ No .env file found, relying on environment variables")
	}
}

// InitDB initializes and returns a PostgreSQL database connection.
func InitDB() *sql.DB {
	InitEnv()

	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	sslmode := os.Getenv("DBSSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open PostgreSQL DB: ", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal("failed to connect to PostgreSQL DB: ", err)
	}

	log.Println("Connected to PostgreSQL DB ✅")
	return db
}

// InitSession initializes the global session store using secure cookies.
func InitSession() {
	InitEnv()

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY is not set in environment")
	}

	Store = sessions.NewCookieStore([]byte(sessionKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
	}
}

// EnsureDataDir creates the data directory with restricted permissions.
func EnsureDataDir() {
	if err := os.MkdirAll("data", 0700); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}
}
