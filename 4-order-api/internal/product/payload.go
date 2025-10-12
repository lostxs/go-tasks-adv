package product

type CreateRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

type UpdateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}
