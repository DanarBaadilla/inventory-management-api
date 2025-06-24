package web

type ProductCreateOrUpdateRequest struct {
	Name       string `json:"name" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required"`
	Stock      int    `json:"stock" validate:"gte=0"`
}
