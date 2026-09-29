

package model

import "time"

type StockMovement struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}