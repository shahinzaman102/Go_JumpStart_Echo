package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/data"
)

// Number is a generic constraint for numeric types (~ allows custom types with same underlying type)
type Number interface {
	~int64 | ~float64
}

// GetTotalBookPrice calculates and returns the total price of all books as JSON.
func GetTotalBookPrice(c echo.Context) error {
	books := data.GetAllBooks()

	// Build a map of ID -> Price
	priceMap := make(map[int]float64)
	for _, b := range books {
		priceMap[b.ID] = b.Price
	}

	total := SumNumbers(priceMap)

	return c.JSON(200, map[string]float64{
		"total_price": total,
	})
}

// SumNumbers adds up all values in a map using generics.
func SumNumbers[K comparable, V Number](m map[K]V) V {
	var total V
	for _, v := range m {
		total += v
	}
	return total
}
