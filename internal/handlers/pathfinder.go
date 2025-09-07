package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shahinzaman102/Go_JumpStart_Echo/internal/pathfinder"
)

// PathfinderResponse represents the result of a pathfinding algorithm.
type PathfinderResponse struct {
	Algorithm     string  `json:"algorithm"`
	PathLength    int     `json:"path_length"`
	ExecutionTime float64 `json:"execution_time_ms"`
}

// Pathfinder runs pathfinding algorithms (Brute Force and BFS)
// on a sample grid and returns their performance as JSON.
// Query parameters: start=x,y and end=x,y
func Pathfinder(c echo.Context) error {
	startParam := c.QueryParam("start")
	endParam := c.QueryParam("end")

	if startParam == "" || endParam == "" {
		return c.String(http.StatusBadRequest, "Missing required params: start, end")
	}

	startCoords := strings.Split(startParam, ",")
	endCoords := strings.Split(endParam, ",")

	if len(startCoords) != 2 || len(endCoords) != 2 {
		return c.String(http.StatusBadRequest, "Invalid coordinates format. Use start=x,y end=x,y")
	}

	startX, err := strconv.Atoi(startCoords[0])
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid start X coordinate")
	}
	startY, err := strconv.Atoi(startCoords[1])
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid start Y coordinate")
	}

	destX, err := strconv.Atoi(endCoords[0])
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid destination X coordinate")
	}
	destY, err := strconv.Atoi(endCoords[1])
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid destination Y coordinate")
	}

	// Sample grid (0 = free cell, 1 = obstacle)
	grid := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 1, 0},
		{0, 0, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 1, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 1, 0, 0},
	}

	rows := len(grid)
	cols := len(grid[0])

	// Validate bounds
	if startX < 0 || startY < 0 || startX >= rows || startY >= cols {
		return c.String(http.StatusBadRequest, "Invalid start coordinates: out of bounds")
	}
	if destX < 0 || destY < 0 || destX >= rows || destY >= cols {
		return c.String(http.StatusBadRequest, "Invalid destination coordinates: out of bounds")
	}

	// Validate obstacles
	if grid[startX][startY] == 1 {
		return c.String(http.StatusBadRequest, "Invalid start coordinates: cannot be on an obstacle")
	}
	if grid[destX][destY] == 1 {
		return c.String(http.StatusBadRequest, "Invalid destination coordinates: cannot be on an obstacle")
	}

	const runs = 1000
	results := []PathfinderResponse{}

	// Brute Force
	var brutePathLen int
	startTime := time.Now()
	for i := 0; i < runs; i++ {
		visited := make([][]bool, len(grid))
		for r := range visited {
			visited[r] = make([]bool, len(grid[0]))
		}
		path := pathfinder.ShortestPathBruteForce(grid, startX, startY, destX, destY, visited)
		if path == math.MaxInt32 {
			path = -1
		}
		brutePathLen = path
	}
	bruteDuration := time.Since(startTime).Seconds() * 1000 / runs

	results = append(results, PathfinderResponse{
		Algorithm:     "brute",
		PathLength:    brutePathLen,
		ExecutionTime: bruteDuration,
	})

	// BFS
	var bfsPathLen int
	startTime = time.Now()
	for i := 0; i < runs; i++ {
		bfsPathLen = pathfinder.ShortestPathBFS(grid, startX, startY, destX, destY)
	}
	bfsDuration := time.Since(startTime).Seconds() * 1000 / runs

	results = append(results, PathfinderResponse{
		Algorithm:     "bfs",
		PathLength:    bfsPathLen,
		ExecutionTime: bfsDuration,
	})

	return c.JSON(http.StatusOK, results)
}
