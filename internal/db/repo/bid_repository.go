package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"tender-managment/internal/model"
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) UpdateBidStatus(bid *model.Bid) error {
	query := `UPDATE bids SET status = $1 WHERE id = $2 AND contractor_id = $3 RETURNING id, contractor_id, tender_id, price, delivery_time, comments, status`

	err := r.db.QueryRow(query, bid.Status, bid.ID, bid.ContractorID).Scan(&bid.ID, &bid.ContractorID, &bid.TenderID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("bid not found or you do not have access to this bid")
		}
		return fmt.Errorf("failed to update bid status: %w", err)
	}

	return nil
}

func (r *BidRepository) GetBidsByContractorID(contractorID int) ([]model.Bid, error) {
	var bids []model.Bid

	query := `SELECT id, contractor_id, tender_id, price, delivery_time, comments, created_at 
			  FROM bids WHERE contractor_id = $1`

	rows, err := r.db.Query(query, contractorID)
	if err != nil {
		return nil, fmt.Errorf("error fetching bids: %w", err)
	}

	for rows.Next() {
		var bid model.Bid
		if err := rows.Scan(&bid.ID, &bid.ContractorID, &bid.TenderID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning bid row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return bids, nil
}

func (r *BidRepository) CreateBid(bid model.Bid) (*model.Bid, error) {
	query := `
		INSERT INTO bids (tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at;
	`
	row := r.db.QueryRow(query, bid.TenderID, bid.ContractorID, bid.Price, bid.DeliveryTime, bid.Comments, bid.Status)
	err := row.Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}
	return &bid, nil
}

func (r *BidRepository) GetBidsByTenderID(tenderID int) ([]model.Bid, error) {
	var bids []model.Bid
	query := `
		SELECT id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
		FROM bids
		WHERE tender_id = $1;
	`
	rows, err := r.db.Query(query, tenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bid model.Bid
		err := rows.Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status, &bid.CreatedAt, &bid.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return bids, nil
}

func (r *BidRepository) GetBidByID(id int) (*model.Bid, error) {
	var bid model.Bid
	query := `
		SELECT id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
		FROM bids	
		WHERE id = $1;
	`
	err := r.db.QueryRow(query, id).Scan(&bid.ID, &bid.TenderID, &bid.ContractorID, &bid.Price, &bid.DeliveryTime, &bid.Comments, &bid.Status, &bid.CreatedAt, &bid.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bid with ID %d: %w", id, err)
	}
	return &bid, nil
}

func (r *BidRepository) DeleteBid(id int) error {
	query := `
		DELETE FROM bids
		WHERE id = $1;
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bid with ID %d: %w", id, err)
	}
	return nil
}

func (r *BidRepository) AwardBid(bidID int) error {
	query := `
		UPDATE bids
		SET status = 'awarded', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`
	_, err := r.db.Exec(query, bidID)
	if err != nil {
		return fmt.Errorf("")
	}
	return nil
}
