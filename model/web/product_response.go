package web

type ProductResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}
