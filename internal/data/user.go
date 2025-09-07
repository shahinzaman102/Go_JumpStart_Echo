package data

import (
	"time"

	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// GetAllUsers fetches all users from the database.
func GetAllUsers() ([]models.User, error) {
	rows, err := db.Query(`SELECT id, username, password, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()

}

// GetUserByID fetches a user by their ID.
func GetUserByID(id int) (*models.User, error) {
	var u models.User
	err := db.QueryRow(`SELECT id, username, password, created_at FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser inserts a new user with hashed password and returns the new user ID.
func CreateUser(username, password string) (int64, error) {
	hashed, err := HashPassword(password)
	if err != nil {
		return 0, err
	}
	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`,
		username, hashed, time.Now())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// HashPassword generates a bcrypt hash of the password.
// Recommended cost for most apps: 10–14 (higher = more secure but slower).
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// UpdateUserByID updates username and/or password for a given user ID.
func UpdateUserByID(id int, username, password string) error {
	if username != "" && password != "" {
		hashed, err := HashPassword(password)
		if err != nil {
			return err
		}
		_, err = db.Exec(`UPDATE users SET username = ?, password = ? WHERE id = ?`, username, hashed, id)
		return err
	}

	if username != "" {
		_, err := db.Exec(`UPDATE users SET username = ? WHERE id = ?`, username, id)
		return err
	}

	if password != "" {
		hashed, err := HashPassword(password)
		if err != nil {
			return err
		}
		_, err = db.Exec(`UPDATE users SET password = ? WHERE id = ?`, hashed, id)
		return err
	}

	return nil
}

// DeleteUserByID removes a user from the database by ID.
func DeleteUserByID(id int) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
