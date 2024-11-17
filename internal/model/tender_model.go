package model

import "time"

type Tender struct {
	ID          int       `json:"id"`
	ClientID    int       `json:"client_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      float64   `json:"budget"`
	Status      string    `json:"status"`
	Attachment  string    `json:"attachment"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTender struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Deadline    string  `json:"deadline" binding:"required"`
	Budget      float64 `json:"budget" binding:"required"`
	Attachment  string  `json:"attachment,omitempty"`
}

type GetTender struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      float64   `json:"budget"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	BidsCount   int       `json:"bids_count"`
}

type UpdateTenderStatusRequest struct {
	Status string `json:"status"`
}
