package product

import "github.com/lib/pq"

type ProductCreate struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" `
	Quantity    uint           `json:"quantity"`
}
type ProductUpdate struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" `
	Quantity    uint           `json:"quantity"`
}
