package web

type CategoryCreateOrUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}
