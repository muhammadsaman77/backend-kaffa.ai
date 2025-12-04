package dto

type CreateProductRequest struct {
	StoreID     string  `json:"store_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"required,min=10,max=500"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	IsAvailable bool    `json:"is_available" binding:"required"`
}
