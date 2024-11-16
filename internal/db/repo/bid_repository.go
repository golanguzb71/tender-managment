package repository

import (
	"database/sql"
	"tender-managment/internal/model"
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{db: db}
}

// CreateBid inserts a new bid into the database and returns the bid ID.
func (r *BidRepository) CreateBid(tenderID, contractorID int, price float64, deliveryTime int, comments string) (int, error) {
	query := `INSERT INTO bids (tender_id, contractor_id, price, delivery_time, comments) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var bidID int
	err := r.db.QueryRow(query, tenderID, contractorID, price, deliveryTime, comments).Scan(&bidID)
	return bidID, err
}

// GetBidsByTenderID retrieves all bids for a given tender ID.
func (r *BidRepository) GetBidsByTenderID(tenderID int) ([]model.Bid, error) {
	query := `SELECT id, contractor_id, price, delivery_time, comments, status 
			  FROM bids WHERE tender_id = $1`
	rows, err := r.db.Query(query, tenderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []model.Bid
	for rows.Next() {
		var bid model.Bid
		if err := rows.Scan(&bid.ID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

// GetBidByID retrieves a single bid by its ID.
func (r *BidRepository) GetBidByID(bidID int) (*model.Bid, error) {
	query := `SELECT id, contractor_id, price, delivery_time, comments, status 
			  FROM bids WHERE id = $1`
	row := r.db.QueryRow(query, bidID)

	var bid model.Bid
	if err := row.Scan(&bid.ID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bid, nil
}

// DeleteBid deletes a bid by its ID.
func (r *BidRepository) DeleteBid(bidID int) error {
	query := `DELETE FROM bids WHERE id = $1`
	_, err := r.db.Exec(query, bidID)
	return err
}

// UpdateBidStatus updates the status of a bid.
func (r *BidRepository) UpdateBidStatus(bidID int, status string) error {
	query := `UPDATE bids SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, bidID)
	return err
}
