package models

type OrderRequest struct {
	AlbumID  int64 `json:"album_id"`
	Quantity int64 `json:"quantity"`
	Customer int64 `json:"customer_id"`
}
