package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// UserResponse defines the JSON output for API clients (hides password)
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// mapUser converts internal User model to UserResponse
func mapUser(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// GetUsers returns a list of all users
func GetUsers(c echo.Context) error {
	users, err := data.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching users"})
	}

	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		userResp := mapUser(u)
		resp[i] = &userResp
	}

	return c.JSON(http.StatusOK, resp)
}

// GetUserByID returns a single user by ID
func GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	user, err := data.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, mapUser(*user))
}

// CreateUser adds a new user to the database
func CreateUser(c echo.Context) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Trim whitespace
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// Validation
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}
	if len(input.Username) < 3 || len(input.Username) > 50 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username must be between 3 and 50 characters"})
	}
	if len(input.Password) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 6 characters"})
	}

	id, err := data.CreateUser(input.Username, input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user"})
	}

	user, _ := data.GetUserByID(int(id))
	return c.JSON(http.StatusCreated, map[string]any{
		"status":  "success",
		"message": "User created successfully",
		"user":    mapUser(*user),
	})
}

// UpdateUser updates username and/or password for a given user
func UpdateUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	if input.Username == "" && input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No fields to update"})
	}
	if input.Username != "" && (len(input.Username) < 3 || len(input.Username) > 50) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username must be between 3 and 50 characters"})
	}
	if input.Password != "" && len(input.Password) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 6 characters"})
	}

	if err := data.UpdateUserByID(id, input.Username, input.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user"})
	}

	updatedUser, err := data.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching updated user"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "success",
		"user":   mapUser(*updatedUser),
	})
}

// DeleteUser removes a user by ID
func DeleteUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := data.DeleteUserByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error deleting user"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "deleted",
		"id":     id,
	})
}
