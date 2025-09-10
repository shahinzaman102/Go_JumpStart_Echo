package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // ensure mysql driver is imported
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	Store *sessions.CookieStore
)

// InitEnv loads .env file so os.Getenv() works
func InitEnv() {
	err := godotenv.Load(".env") // adjust path if needed
	if err != nil {
		log.Println("⚠️ No .env file found, relying on environment variables")
	}
}

// InitDB initializes and returns a MySQL database connection.
func InitDB() *sql.DB {
	InitEnv() // Load env variables

	// Config (Local DB)
	// -----------------
	// user := os.Getenv("DBUSER")
	// pass := os.Getenv("DBPASS")
	// name := os.Getenv("DBNAME")

	// DSN (Local DB)
	// --------------
	// dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true&multiStatements=true", user, pass, name)

	// Config (freesqldatabase.com)
	// ----------------------------
	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	name := os.Getenv("DBNAME")

	// DSN (freesqldatabase.com)
	// -------------------------
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to open DB: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}

	log.Println("Connected to DB ✅")
	return db
}

// InitSession initializes the global session store using secure cookies.
func InitSession() {
	InitEnv() // ensure .env is loaded

	// Get session key from env
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY is not set in environment")
	}

	Store = sessions.NewCookieStore([]byte(sessionKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Secure:   false, // set true if using HTTPS
	}
}

// EnsureDataDir creates the data directory with restricted permissions (owner-only).
func EnsureDataDir() {
	if err := os.MkdirAll("data", 0700); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}
}
