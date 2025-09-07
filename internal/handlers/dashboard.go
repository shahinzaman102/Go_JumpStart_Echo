package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/models"
)

// Dashboard shows tasks if the user is logged in; otherwise, shows 401 + login page.
func Dashboard(c echo.Context) error {
	if !isAuthenticated(c) {
		session, _ := store.Get(c.Request(), "session")
		session.Values["redirect_after_login"] = "/dashboard"
		session.Save(c.Request(), c.Response())

		c.Response().WriteHeader(http.StatusUnauthorized)
		c.Response().Header().Set("Content-Type", "text/html")

		tmpl := template.Must(template.ParseFiles("templates/login_required.html"))
		return tmpl.Execute(c.Response(), map[string]string{"Resource": "dashboard"})
	}

	// Example due dates
	dueBuildApp := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	dueLearnGo := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	todos := []models.Todo{
		{Title: "Learn Go", Done: true, Progress: 100, Due: &dueLearnGo},
		{Title: "Build the 2nd App", Done: false, Progress: 40, Due: &dueBuildApp},
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	return tmpl.Execute(c.Response(), todos)
}

// isAuthenticated checks session for login status
func isAuthenticated(c echo.Context) bool {
	session, _ := store.Get(c.Request(), "session")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}
