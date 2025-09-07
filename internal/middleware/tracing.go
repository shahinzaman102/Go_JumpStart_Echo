package middleware

import (
	"fmt"
	"runtime/trace"
	"time"

	"github.com/labstack/echo/v4"
)

// Tracing is an Echo middleware that tracks request details and execution time using Go's runtime trace.
func Tracing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start a new trace task for this request
		ctx, task := trace.NewTask(c.Request().Context(), fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path))
		defer task.End()

		// Log request details
		trace.Log(ctx, "method", c.Request().Method)
		trace.Log(ctx, "path", c.Request().URL.Path)

		// Measure request duration
		start := time.Now()
		c.SetRequest(c.Request().WithContext(ctx)) // pass tracing context to downstream handlers
		err := next(c)
		trace.Log(ctx, "duration_ms", fmt.Sprintf("%.2f", time.Since(start).Seconds()*1000))

		return err
	}
}
