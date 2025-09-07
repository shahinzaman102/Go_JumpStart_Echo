package handlers

import (
	"io/fs"
	"net/http"

	assets "github.com/shahinzaman102/Go_JumpStart_Echo"

	"github.com/labstack/echo/v4"
)

// ServeStaticEcho registers embedded static files for Echo.
// Example: /static/style.css serves embedded static/style.css
func ServeStaticEcho(e *echo.Echo) {
	staticFiles, _ := fs.Sub(assets.StaticFS, "static")
	fsHandler := http.FileServer(http.FS(staticFiles))

	// Use Echo's GET handler with a wildcard
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fsHandler)))
}
