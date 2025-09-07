package handlers

import (
	"html/template"

	assets "github.com/shahinzaman102/Go_JumpStart_Echo"

	"github.com/labstack/echo/v4"
)

// TestUI serves the test_ui.html template via Echo
func TestUI(c echo.Context) error {
	tmpl := template.Must(template.ParseFS(assets.Templates, "templates/test_ui.html"))

	// Render template to response
	return tmpl.ExecuteTemplate(c.Response().Writer, "test_ui.html", nil)
}
