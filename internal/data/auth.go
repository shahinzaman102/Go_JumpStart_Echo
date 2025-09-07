package data

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// AuthRepo provides authentication methods using a SQL database.
type AuthRepo struct {
	DB *sql.DB
}

// NewAuthRepo creates a new AuthRepo with a given DB connection.
func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

// VerifyUser checks if the username/password combination is valid.
func (repo *AuthRepo) VerifyUser(username, password string) (bool, error) {
	var hash string
	err := repo.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hash)
	if err != nil {
		return false, err
	}
	return CheckPasswordHash(password, hash), nil
}

// CheckPasswordHash compares a plain password with a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserID retrieves the ID of a user by username.
func (repo *AuthRepo) GetUserID(username string) (int64, error) {
	var id int64
	err := repo.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
