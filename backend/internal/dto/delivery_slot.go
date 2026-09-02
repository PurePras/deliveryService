package dto

type DeliverySlotInput struct {
	Label     string `json:"label"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsActive  *bool  `json:"is_active"`
}
