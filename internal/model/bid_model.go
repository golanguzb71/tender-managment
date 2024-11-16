package model

import "time"

type Bid struct {
	ID           string    `json:"id" bson:"_id"`
	Price        float64   `json:"price"`
	DeliveryTime int       `json:"delivery_time"`
	Comments     string    `json:"comments"`
	ContractorID int       `json:"contractor_id"`
	TenderID     int       `json:"tender_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateBid struct {
	Price        float64 `json:"price"`
	DeliveryTime int     `json:"delivery_time"`
	Comments     string  `json:"comments"`
}

const (
	BidStatusPending  = "pending"
	BidStatusAwarded  = "awarded"
	BidStatusRejected = "rejected"
)
