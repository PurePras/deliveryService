package dto

type CategoryInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}
