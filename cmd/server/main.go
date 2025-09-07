package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime/trace"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/config"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/db"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/handlers"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/routes"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// --- Start runtime tracing ---
	traceFile, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		traceFile.Close()
		log.Println("trace file closed")
	}()
	if err := trace.Start(traceFile); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer func() {
		trace.Stop()
		log.Println("trace stopped")
	}()
	log.Println("Go trace started")

	// --- Ensure data directory exists ---
	config.EnsureDataDir()

	// --- Initialize DB, sessions, cache ---
	conn := config.InitDB()
	defer func() {
		conn.Close()
		log.Println("database connection closed")
	}()
	config.InitSession()
	data.InitDBConnection(conn)
	data.InitCache()
	authRepo := data.NewAuthRepo(conn)
	handlers.Init(config.Store, authRepo)

	// --- Preload wiki templates ---
	if err := handlers.LoadWikiTemplates(); err != nil {
		log.Fatalf("failed to load wiki templates: %v", err)
	}

	// --- Execute schema ---
	if err := db.ExecuteSchema(conn, "schema.sql"); err != nil {
		log.Fatalf("schema execution failed: %v", err)
	}
	log.Println("database schema executed successfully")

	// --- Initialize Echo ---
	e := echo.New()

	// Hide Echo banner & port
	e.HideBanner = true
	e.HidePort = true

	// Disable Echo default logger completely
	e.Logger.SetOutput(io.Discard)

	// Only use Recover middleware
	e.Use(middleware.Recover())

	// --- Minimal logging middleware (skip static files) ---
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			// Skip logging static file requests
			if len(c.Path()) < 8 || c.Path()[:8] != "/static/" {
				status := c.Response().Status
				log.Printf("%s %s → %d", c.Request().Method, c.Request().RequestURI, status)
			}
			return err
		}
	})

	// --- Register routes ---
	routes.Register(e)

	// --- Start pprof server in background ---
	go func() {
		log.Println("pprof server listening on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof server error: %v", err)
		}
	}()

	// --- Start Echo server ---
	log.Println("server running at http://localhost:8080")
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("main server error: %v", err)
	}

	log.Println("server stopped")
}
