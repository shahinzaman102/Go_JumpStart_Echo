package handlers

import (
	"html/template"
	"net/http"

	assets "github.com/shahinzaman102/Go_JumpStart_Echo"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var (
	store    *sessions.CookieStore
	authRepo *data.AuthRepo
)

// Init initializes the session store and auth repository.
func Init(s *sessions.CookieStore, repo *data.AuthRepo) {
	store = s
	authRepo = repo
}

// Login handles user login: verifies credentials and sets session values.
func Login(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")

	username := c.FormValue("username")
	password := c.FormValue("password")

	ok, err := authRepo.VerifyUser(username, password)
	if err != nil || !ok {
		tmpl := template.Must(template.ParseFS(assets.Templates, "templates/unauthorized.html"))
		c.Response().WriteHeader(http.StatusUnauthorized)
		return tmpl.Execute(c.Response(), nil)
	}

	userID, err := authRepo.GetUserID(username)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to fetch user ID")
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = userID

	// Determine redirect path
	redirectPath := c.FormValue("redirect")
	if redirectPath == "" {
		redirectPath = c.QueryParam("redirect")
	}
	if redirectPath == "" {
		if v, ok := session.Values["redirect_after_login"].(string); ok && v != "" {
			redirectPath = v
			delete(session.Values, "redirect_after_login")
		}
	}
	if redirectPath == "" {
		redirectPath = "/dashboard"
	}

	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save session")
	}

	return c.Redirect(http.StatusSeeOther, redirectPath)
}

// LoginForm renders the login page with an optional redirect.
func LoginForm(c echo.Context) error {
	redirect := c.QueryParam("redirect")
	tmpl := template.Must(template.ParseFS(assets.Templates, "templates/login.html"))
	return tmpl.Execute(c.Response(), map[string]string{
		"Redirect": redirect,
	})
}

// Logout clears the session and redirects to the home page.
func Logout(c echo.Context) error {
	session, _ := store.Get(c.Request(), "session")

	// Clear session completely
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)

	if err := session.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to clear session")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
