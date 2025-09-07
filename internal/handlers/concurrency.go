package handlers

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	Atomic_Counters "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Atomic_Counters"
	Buffered_Channels "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Buffered_Channels"
	Channels_unbuffered "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Channels_unbuffered"
	Cond_syncCond "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Cond_syncCond"
	Context_Cancellation "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Context_Cancellation"
	Goroutines_WaitGroup "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Goroutines_WaitGroup"
	Mutex "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Mutex"
	Pool_Once_Map "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Pool_Once_Map"
	RWMutex "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/RWMutex"
	Worker_Pool "github.com/shahinzaman102/Go_JumpStart_Echo/internal/concurrency/Worker_Pool"
)

// captureOutput redirects stdout to a buffer and returns it as string
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// helper to write captured output as plain text
func writeDemoOutput(c echo.Context, runFunc func()) error {
	return c.String(http.StatusOK, captureOutput(runFunc))
}

// ==== Handlers ====
func GoroutinesWaitGroupHandler(c echo.Context) error {
	return writeDemoOutput(c, Goroutines_WaitGroup.Run)
}

func ChannelsUnbufferedHandler(c echo.Context) error {
	return writeDemoOutput(c, Channels_unbuffered.Run)
}

func BufferedChannelsHandler(c echo.Context) error {
	return writeDemoOutput(c, Buffered_Channels.Run)
}

func MutexHandler(c echo.Context) error {
	return writeDemoOutput(c, Mutex.Run)
}

func RWMutexHandler(c echo.Context) error {
	return writeDemoOutput(c, RWMutex.Run)
}

func WorkerPoolHandler(c echo.Context) error {
	return writeDemoOutput(c, Worker_Pool.Run)
}

func AtomicCountersHandler(c echo.Context) error {
	return writeDemoOutput(c, Atomic_Counters.Run)
}

func CondSyncCondHandler(c echo.Context) error {
	return writeDemoOutput(c, Cond_syncCond.Run)
}

func PoolOnceMapHandler(c echo.Context) error {
	return writeDemoOutput(c, Pool_Once_Map.Run)
}

func ContextCancellationHandler(c echo.Context) error {
	return writeDemoOutput(c, Context_Cancellation.Run)
}
