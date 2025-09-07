package routes

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware" // alias Echo middleware

	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/handlers"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/middleware"
)

// Register registers all routes with Echo
func Register(e *echo.Echo) {
	// --- Middleware ---
	e.Use(echomw.Logger())  // Echo logger
	e.Use(echomw.Recover()) // Echo recover
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	e.Use(middleware.Tracing) // your custom tracing middleware

	// --- App Home ---
	e.GET("/", handlers.TestUI)

	// --- Form ---
	e.GET("/form", handlers.Form)
	e.POST("/form", handlers.Form)

	// --- Static Files ---
	handlers.ServeStaticEcho(e)

	// --- Authentication Flow ---
	e.GET("/login", handlers.LoginForm)
	e.POST("/login", handlers.Login)
	e.GET("/logout", handlers.Logout)

	// --- Dashboard ---
	e.GET("/dashboard", handlers.Dashboard)

	// --- Users API ---
	users := e.Group("/users")
	users.GET("", handlers.GetUsers)
	users.POST("", handlers.CreateUser)
	users.GET("/:id", handlers.GetUserByID)
	users.PUT("/:id", handlers.UpdateUser)
	users.DELETE("/:id", handlers.DeleteUser)

	// --- Books API ---
	books := e.Group("/books")
	books.GET("", handlers.GetBooks)
	books.POST("", handlers.PostBook)
	books.GET("/total", handlers.GetTotalBookPrice)
	books.GET("/:id", handlers.GetBookByID)
	books.PUT("/:id", handlers.UpdateBook)
	books.DELETE("/:id", handlers.DeleteBook)

	// --- Albums API ---
	albums := e.Group("/albums")
	albums.GET("", handlers.GetAllAlbums)
	albums.POST("", handlers.CreateAlbum)
	albums.GET("/artist/:name", handlers.GetAlbumsByArtist)
	albums.GET("/timeout", handlers.QueryWithTimeout)
	albums.GET("/:id/can-purchase", handlers.CanPurchaseAlbum)
	albums.GET("/:id", handlers.GetAlbumByID)

	// --- Orders API ---
	orders := e.Group("/orders")
	orders.GET("", handlers.GetOrdersByUser)
	orders.POST("", handlers.CreateOrderByUser)

	// --- Misc Handlers ---
	e.GET("/customer-name", handlers.GetCustomerName)
	e.GET("/admin/multi-query", handlers.HandleMultipleResultSets)

	// --- Wiki Pages ---
	e.GET("/view", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/view/FrontPage")
	})
	e.GET("/view/:title", handlers.ViewWiki)
	e.GET("/edit/:title", handlers.EditWiki)
	e.POST("/save/:title", handlers.SaveWiki)

	// --- JSON Utilities ---
	e.POST("/json/encode", handlers.JsonEncode)
	e.POST("/json/decode", handlers.JsonDecode)

	// --- WebSocket ---
	e.GET("/ws", handlers.Echo) // WebSocket upgrade
	e.GET("/websockets", handlers.WebsocketPage)

	// --- Concurrency ---
	e.GET("/concurrency/goroutines_waitgroup", handlers.GoroutinesWaitGroupHandler)
	e.GET("/concurrency/channels_unbuffered", handlers.ChannelsUnbufferedHandler)
	e.GET("/concurrency/buffered_channels", handlers.BufferedChannelsHandler)
	e.GET("/concurrency/mutex", handlers.MutexHandler)
	e.GET("/concurrency/rwmutex", handlers.RWMutexHandler)
	e.GET("/concurrency/worker_pool", handlers.WorkerPoolHandler)
	e.GET("/concurrency/atomic_counters", handlers.AtomicCountersHandler)
	e.GET("/concurrency/cond_synccond", handlers.CondSyncCondHandler)
	e.GET("/concurrency/pool_once_map", handlers.PoolOnceMapHandler)
	e.GET("/concurrency/context_cancellation", handlers.ContextCancellationHandler)

	e.GET("/go-basics", handlers.GoBasics)
	e.GET("/pathfinder", handlers.Pathfinder)
	e.GET("/runtime-errors", handlers.RuntimeErrorsHandler)
}
