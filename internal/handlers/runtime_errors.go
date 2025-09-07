package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RuntimeErrorsHandler provides examples of common Go runtime errors.
// Trigger using query params, e.g. /runtime-errors?type=divide
func RuntimeErrorsHandler(c echo.Context) error {
	errType := c.QueryParam("type")

	defer func() {
		if rec := recover(); rec != nil {
			_ = c.String(http.StatusInternalServerError, fmt.Sprintf("Recovered from panic: %v", rec))
		}
	}()

	switch errType {
	case "divide":
		// Division by zero
		x := 10
		y := 0
		_ = x / y
	case "nilptr":
		// Nil pointer dereference
		var p *int
		_ = *p
	case "outofbounds":
		// Indexing outside array/slice
		arr := []int{1, 2, 3}
		_ = arr[10]
	case "typeassert":
		// Invalid type assertion
		var i any = "hello"
		_ = i.(int)
	default:
		return c.String(http.StatusBadRequest, "Invalid type. Use one of: divide, nilptr, outofbounds, typeassert\n")
	}

	// If panic didn’t happen (should not reach here normally)
	return c.String(http.StatusOK, "No error occurred (unexpected).")
}
