package web

import "time"

type StockMovementResponse struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Product   string    `json:"product"`
	UserID    int       `json:"user_id"`
	User      string    `json:"user"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
