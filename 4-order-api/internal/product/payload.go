package product

type CreateRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       float64  `json:"price" validate:"required"`
}

type UpdateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       float64  `json:"price"`
}
