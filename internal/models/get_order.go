package models

import "time"

type GetOrder struct {
	ID       int64     `json:"id"`
	AlbumID  int64     `json:"album_id"`
	Customer int64     `json:"customer_id"`
	Quantity int64     `json:"quantity"`
	Date     time.Time `json:"date"`
}
