package web

type StockMovementCreateRequest struct {
	ProductID int    `json:"product_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=in out"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
	Note      string `json:"note"`
}
