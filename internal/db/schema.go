package db

import (
	"database/sql"
	"fmt"
	"os"
)

// ExecuteSchema reads a SQL file and executes its statements on the provided DB.
func ExecuteSchema(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute all statements in the file
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}
	return nil
}
