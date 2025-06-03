package dtos

type ProductDto struct {
	ID          int     `json:"id"`          // ID is the primary key, not required for creation
	Title       string  `json:"title"`       // Title is required for product creation
	Description *string `json:"description"` // Review is optional, can be nil
	Image       *string `json:"image"`       // Image is optional, can be nil
	Price       float64 `json:"price"`       // Price is required for product creation
}
