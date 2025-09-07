package handlers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GoBasics demonstrates core Go features through HTTP output.
func GoBasics(c echo.Context) error {
	var output string

	appendLine := func(format string, a ...any) {
		output += fmt.Sprintf(format+"\n", a...)
	}

	// --- Type Conversion ---
	priceStr := "199"
	price, err := strconv.Atoi(priceStr) // string → int
	if err != nil {
		appendLine("Error converting price: %v", err)
	} else {
		appendLine("Converted string '%s' → int: %d", priceStr, price)
	}

	discount := 15
	finalPrice := float64(price) * (1 - float64(discount)/100)
	appendLine("Final Price after %d%% discount: %.2f", discount, finalPrice)

	rating := 4.8
	appendLine("User rating %.1f stored as int stars: %d", rating, int(rating))

	arr := []string{"Go", "Python", "Rust"}
	uidx := uint(2)
	appendLine("Accessing array[%d] → %s", uidx, arr[uidx])

	appendLine("")

	// --- Switch ---
	day := 2
	switch day {
	case 1:
		appendLine("Switch Case: Friday")
	case 2:
		appendLine("Switch Case: Saturday")
	case 3:
		appendLine("Switch Case: Sunday")
	default:
		appendLine("Switch Case: Unknown Day")
	}

	appendLine("")

	// --- iota ---
	const (
		First = iota
		Second
		Third
	)
	appendLine("Basic iota: First=%d, Second=%d, Third=%d", First, Second, Third)

	const (
		FlagRead = 1 << iota
		FlagWrite
		FlagExecute
	)
	appendLine("Bit Flags: Read=%d, Write=%d, Execute=%d", FlagRead, FlagWrite, FlagExecute)

	appendLine("")

	// --- Command-line arguments ---
	appendLine("Command-line arguments:")
	for i, arg := range os.Args {
		appendLine("Arg %d: %s", i, arg)
	}
	if len(os.Args) <= 1 {
		appendLine("\nTip: run with args, e.g.: go run ./cmd/server arg1 arg2 arg3")
	}

	appendLine("")

	// --- Panic and Recover ---
	processPayment := func(amount float64) {
		defer func() {
			if r := recover(); r != nil {
				appendLine("Recovered from failure: %v", r)
			}
		}()

		appendLine("Connecting to payment gateway...")
		if amount <= 0 {
			panic("invalid payment amount")
		}

		_, err := os.Open("credentials.txt")
		if err != nil {
			panic("missing payment credentials file")
		}

		appendLine("Payment processed successfully: %.2f", amount)
	}

	appendLine("\n=== Payment Demo ===")
	processPayment(0)   // triggers panic, recovered
	processPayment(100) // may panic if file missing

	appendLine("")

	// --- Exit Codes (simulated for safety in HTTP) ---
	if exitVal := c.QueryParam("exit"); exitVal != "" {
		code, err := strconv.Atoi(exitVal)
		if err != nil {
			return c.String(400, "Invalid exit code")
		}

		var msg string
		switch code {
		case 0:
			msg = "Exit code 0: Success"
		case 1:
			msg = "Exit code 1: Generic error"
		case 2:
			msg = "Exit code 2: Invalid input"
		case 5:
			msg = "Exit code 5: Config missing"
		default:
			msg = fmt.Sprintf("Exit code %d: Unknown reason", code)
		}

		appendLine(msg)
		appendLine("Server would exit with code %d (skipped for safety)", code)
	}

	appendLine("\n💡 Safe URLs to test exit codes:")
	appendLine("http://localhost:8080/go-basics?exit=0")
	appendLine("http://localhost:8080/go-basics?exit=1")
	appendLine("http://localhost:8080/go-basics?exit=2")
	appendLine("http://localhost:8080/go-basics?exit=5")
	appendLine("http://localhost:8080/go-basics?exit=999")

	appendLine("\nGo Basics Demo Complete.")

	return c.String(200, output)
}
