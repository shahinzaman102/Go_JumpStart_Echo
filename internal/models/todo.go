package models

import "time"

type Todo struct {
	Title    string
	Done     bool
	Progress int        // 0-100%
	Due      *time.Time // nil if no due date
}
