package model

type Gig struct {
	Id           uint64  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	DeliveryDays int32   `json:"delivery_days"`
	Revisions    int32   `json:"revisions"`
	FreelancerId uint64  `json:"freelancer_id"`
	Images       string  `json:"images"`
}
