package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// upgrader upgrades HTTP requests to WebSocket connections
var upgrader = websocket.Upgrader{
	// Allow all origins for demo purposes (adjust for production)
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Echo is a WebSocket handler that echoes messages back to the client
func Echo(c echo.Context) error {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return err
	}
	defer conn.Close()

	// Continuously read messages from the client and echo them back
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
	return nil
}

// WebsocketPage serves the HTML page for the WebSocket frontend
func WebsocketPage(c echo.Context) error {
	tmpl, err := template.ParseFiles("templates/websockets.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Template parse error: "+err.Error())
	}

	return tmpl.Execute(c.Response(), nil)
}
