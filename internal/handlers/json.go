package handlers

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// Encode: Take input fields -> struct -> return JSON
func JsonEncode(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.String(400, fmt.Sprintf("Invalid input: %v", err))
	}

	// Validate required fields
	if u.Username == "" || u.Password == "" {
		return c.String(400, "Username and Password are required")
	}

	u.CreatedAt = time.Now() // auto-assign timestamp

	return c.JSON(200, u)
}

// Decode: Take input fields -> struct -> return text summary
func JsonDecode(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.String(400, fmt.Sprintf("Invalid input: %v", err))
	}

	// Validate required fields
	if u.Username == "" || u.Password == "" {
		return c.String(400, "Username and Password are required")
	}

	return c.String(200, fmt.Sprintf("Received user: %+v", u))
}
