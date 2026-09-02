package dto

type ProductInput struct {
	CategoryID    string  `json:"category_id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   *string `json:"description"`
	Unit          string  `json:"unit"`
	Price         string  `json:"price"`
	StockQuantity string  `json:"stock_quantity"`
	ImageURL      *string `json:"image_url"`
	IsAvailable   *bool   `json:"is_available"`
}
