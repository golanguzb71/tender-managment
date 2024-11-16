package service

import (
	"errors"
	"fmt"
	"net/http"
	repository "tender-managment/internal/db/repo"
	"tender-managment/internal/model"
)

type BidService struct {
	bidRepo        repository.BidRepository
	tenderRepo     repository.TenderRepository
	contractorRepo repository.UserRepository
}

func NewBidService(bidRepo repository.BidRepository, tenderRepo repository.TenderRepository, contractorRepo repository.UserRepository) *BidService {
	return &BidService{
		bidRepo:        bidRepo,
		tenderRepo:     tenderRepo,
		contractorRepo: contractorRepo,
	}
}

func (s *BidService) GetBidByID(contractorID int, bidID int) (*model.Bid, error) {
	// Fetch the bid from the repository
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return nil, fmt.Errorf("bid not found: %w", err)
	}

	// Check if the bid belongs to the contractor
	if bid.ContractorID != contractorID {
		return nil, errors.New("you do not have access to this bid")
	}

	return bid, nil
}

// GetBidsByContractor retrieves all bids placed by a specific contractor
func (s *BidService) GetBidsByContractor(contractorID int) ([]model.Bid, error) {
	// Fetch the bids for the contractor
	bids, err := s.bidRepo.GetBidsByContractorID(contractorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bids: %w", err)
	}

	return bids, nil
}

func (s *BidService) CreateBid(contractorID int, tenderID int, bid model.CreateBid) (*model.Bid, int, error) {
	// Fetch the tender to ensure it exists and is open
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Tender not found")
	}

	if tender.Status != "open" {
		return nil, http.StatusBadRequest, fmt.Errorf("Tender is not open for bids")
	}

	// Validate the bid data
	if bid.Price <= 0 || bid.DeliveryTime <= 0 || bid.Comments == "" {
		return nil, http.StatusBadRequest, errors.New("invalid bid data")
	}
	var newBid model.Bid
	newBid.Price = bid.Price
	newBid.Comments = bid.Comments
	newBid.DeliveryTime = bid.DeliveryTime
	newBid.ContractorID = contractorID
	newBid.TenderID = tenderID
	newBid.Status = model.BidStatusPending
	createdBid, err := s.bidRepo.CreateBid(newBid)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to create bid: %w", err)
	}

	return createdBid, http.StatusCreated, nil
}

// GetBidsByTenderID retrieves all bids for a specific tender
func (s *BidService) GetBidsByTenderID(tenderID, userId int) ([]model.Bid, error) {
	// Fetch the tender to ensure it exists
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil || tender.ClientID != userId {
		return nil, fmt.Errorf("Tender not found or access denied")
	}
	fmt.Println(tender)
	//// Check if the user is allowed to view the bids (either client or contractor)
	//if !utils.HasPermission(tender.ClientID) {
	//	return nil, errors.New("access denied")
	//}

	// Get the bids for the tender
	bids, err := s.bidRepo.GetBidsByTenderID(tenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bids: %w", err)
	}

	return bids, nil
}

// DeleteBid handles the logic for deleting a bid
func (s *BidService) DeleteBid(contractorID int, bidID int) error {
	// Fetch the bid to ensure it exists
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", err)
	}

	// Check if the contractor owns the bid
	if bid.ContractorID != contractorID {
		return errors.New("you do not own this bid")
	}

	// Delete the bid
	err = s.bidRepo.DeleteBid(bidID)
	if err != nil {
		return fmt.Errorf("failed to delete bid: %w", err)
	}

	return nil
}

// UpdateBidStatus updates the status of a specific bid for the contractor.
func (s *BidService) UpdateBidStatus(contractorID int, bidID int, newStatus string) error {
	// Validate the status value (ensure it's one of the allowed values)
	validStatuses := []string{model.BidStatusAwarded, model.BidStatusRejected, model.BidStatusRejected}
	if !contains(validStatuses, newStatus) {
		return errors.New("invalid status")
	}

	// Fetch the bid to ensure it exists
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", err)
	}

	// Check if the bid belongs to the contractor
	if bid.ContractorID != contractorID {
		return errors.New("you do not have access to this bid")
	}

	// Update the status of the bid
	bid.Status = newStatus
	err = s.bidRepo.UpdateBidStatus(bid)
	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}

	return nil
}

// Helper function to check if a value exists in a slice
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// AwardBid handles the logic for awarding a bid
func (s *BidService) AwardBid(clientID int, tenderID int, bidID int) error {
	// Fetch the tender to ensure it exists and is owned by the client
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil || tender.ClientID != clientID {
		return fmt.Errorf("Tender not found or access denied")
	}

	// Fetch the bid to ensure it exists
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", err)
	}
	fmt.Println(bid)
	err = s.bidRepo.AwardBid(bidID)
	if err != nil {
		return fmt.Errorf("failed to award bid: %w", err)
	}

	// Update the tender status to awarded
	err = s.tenderRepo.UpdateTenderStatus(tenderID, "awarded")
	if err != nil {
		return fmt.Errorf("failed to update tender status: %w", err)
	}

	return nil
}
