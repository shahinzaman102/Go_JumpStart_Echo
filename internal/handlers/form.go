package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	assets "github.com/shahinzaman102/Go_JumpStart_Echo"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

func Form(c echo.Context) error {
	// Load form.html from embedded templates
	tmpl := template.Must(template.ParseFS(assets.Templates, "templates/form.html"))

	if c.Request().Method == http.MethodPost {
		// Parse submitted form fields
		if err := c.Request().ParseForm(); err != nil {
			log.Println("Form parse error:", err)
			return err
		}

		email := c.Request().FormValue("email")
		subject := c.Request().FormValue("subject")
		message := c.Request().FormValue("message")

		log.Println("Form received:", email, subject, message)

		data := models.FormResponse{
			Success: true,
			Email:   email,
			Subject: subject,
			Message: message,
		}

		return tmpl.Execute(c.Response(), data)
	}

	// Render empty form on GET
	return tmpl.Execute(c.Response(), nil)
}
