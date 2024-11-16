package service

import (
	"errors"
	repository "tender-managment/internal/db/repo"
	"tender-managment/internal/model"
)

type BidService struct {
	bidRepo *repository.BidRepository
}

func NewBidService(bidRepo *repository.BidRepository) *BidService {
	return &BidService{bidRepo: bidRepo}
}

// CreateBid handles the logic of creating a new bid.
func (s *BidService) CreateBid(tenderID, contractorID int, price float64, deliveryTime int, comments string) (int, error) {
	if tenderID <= 0 || contractorID <= 0 || price <= 0 || deliveryTime <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	// Call repository to create the bid
	bidID, err := s.bidRepo.CreateBid(tenderID, contractorID, price, deliveryTime, comments)
	if err != nil {
		return 0, err
	}

	return bidID, nil
}

// GetBidsByTenderID retrieves all bids for a given tender.
func (s *BidService) GetBidsByTenderID(tenderID int) ([]model.Bid, error) {
	if tenderID <= 0 {
		return nil, errors.New("invalid tender ID")
	}

	bids, err := s.bidRepo.GetBidsByTenderID(tenderID)
	if err != nil {
		return nil, err
	}

	return bids, nil
}

// GetBidByID retrieves a single bid by its ID.
func (s *BidService) GetBidByID(bidID int) (*model.Bid, error) {
	if bidID <= 0 {
		return nil, errors.New("invalid bid ID")
	}

	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return nil, err
	}

	return bid, nil
}

// DeleteBid deletes a bid by its ID.
func (s *BidService) DeleteBid(bidID int) error {
	if bidID <= 0 {
		return errors.New("invalid bid ID")
	}

	return s.bidRepo.DeleteBid(bidID)
}

// UpdateBidStatus updates the status of a bid.
func (s *BidService) UpdateBidStatus(bidID int, status string) error {
	if bidID <= 0 || status == "" {
		return errors.New("invalid parameters")
	}

	return s.bidRepo.UpdateBidStatus(bidID, status)
}
