package dto

type CreateProductRequest struct {
	StoreID     string  `form:"store_id" binding:"required"`
	Name        string  `form:"name" binding:"required,min=2,max=100"`
	Description string  `form:"description" binding:"required,min=10,max=500"`
	Price       float64 `form:"price" binding:"required,gt=0"`
	IsAvailable bool    `form:"is_available" binding:"required"`
}

type GetAllProductResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	ImagePath   string  `json:"image_path"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type GetProductDetailsResponse struct {
	Id          string  `json:"id"`
	StoreID     string  `json:"store_id"`
	ImageID     string  `json:"image_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	ImagePath   string  `json:"image_path"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
