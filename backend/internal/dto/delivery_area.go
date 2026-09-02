package dto

type DeliveryAreaInput struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Pincode  string `json:"pincode"`
	IsActive *bool  `json:"is_active"`
}
