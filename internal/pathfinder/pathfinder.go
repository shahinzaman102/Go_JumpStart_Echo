package pathfinder

import (
	"math"
)

// ShortestPathBruteForce finds the shortest path using DFS (backtracking).
// Returns math.MaxInt32 if no path exists.
func ShortestPathBruteForce(grid [][]int, x, y, destX, destY int, visited [][]bool) int {
	rows := len(grid)
	cols := len(grid[0])

	// Out of bounds or blocked/visited cell
	if x < 0 || y < 0 || x >= rows || y >= cols {
		return math.MaxInt32
	}

	// Blocked or visited
	if grid[x][y] == 1 || visited[x][y] {
		return math.MaxInt32
	}
	// math.MaxInt32 = largest 32-bit int → used here as a marker for an unreachable path.

	// Destination reached
	if x == destX && y == destY {
		return 0
	}

	// Mark visited
	visited[x][y] = true

	// Explore 4 directions: left, right, up, down
	left := ShortestPathBruteForce(grid, x, y-1, destX, destY, visited)
	right := ShortestPathBruteForce(grid, x, y+1, destX, destY, visited)
	up := ShortestPathBruteForce(grid, x-1, y, destX, destY, visited)
	down := ShortestPathBruteForce(grid, x+1, y, destX, destY, visited)

	// Unmark current cell for backtracking
	visited[x][y] = false

	minPath := min(down, up, right, left)
	if minPath == math.MaxInt32 {
		return math.MaxInt32
	}
	return 1 + minPath
}

// min returns the smallest value among the given integers.
func min(nums ...int) int {
	minVal := nums[0]
	for _, v := range nums {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

// ShortestPathBFS finds the shortest path using BFS (level-order traversal).
// Returns -1 if no path exists.
func ShortestPathBFS(grid [][]int, startX, startY, destX, destY int) int {
	rows := len(grid)
	cols := len(grid[0])

	type cell struct {
		x, y, dist int
	}

	queue := []cell{{startX, startY, 0}}

	// Track visited cells using a map keyed by coordinates
	visited := make(map[[2]int]bool)
	visited[[2]int{startX, startY}] = true

	// 4 possible directions: left, right, up, down
	directions := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// Destination reached
		if curr.x == destX && curr.y == destY {
			return curr.dist
		}

		for _, dir := range directions {
			nx := curr.x + dir[0] // nx → neighbor row
			ny := curr.y + dir[1] // ny → neighbor column
			if nx >= 0 && ny >= 0 && nx < rows && ny < cols &&
				grid[nx][ny] == 0 && !visited[[2]int{nx, ny}] {
				visited[[2]int{nx, ny}] = true
				queue = append(queue, cell{nx, ny, curr.dist + 1})
			}
		}
	}

	// No path found
	return -1
}
