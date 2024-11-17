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
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return nil, fmt.Errorf("bid not found: %w", err)
	}

	if bid.ContractorID != contractorID {
		return nil, errors.New("you do not have access to this bid")
	}

	return bid, nil
}

func (s *BidService) GetBidsByContractor(contractorID int) ([]model.Bid, error) {
	bids, err := s.bidRepo.GetBidsByContractorID(contractorID)
	if err != nil {
		fmt.Println("error bo'ldi ", err.Error())
		return nil, fmt.Errorf("failed to fetch bids: %w", err)
	}

	return bids, nil
}

func (s *BidService) CreateBid(contractorID int, tenderID int, bid model.CreateBid) (*model.Bid, int, error) {
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("Tender not found")
	}

	if tender.Status != "open" {
		return nil, http.StatusBadRequest, fmt.Errorf("Tender is not open for bids")
	}

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

func (s *BidService) GetBidsByTenderID(tenderID, userId int) ([]model.Bid, error) {
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil || tender.ClientID != userId {
		return nil, fmt.Errorf("Tender not found or access denied")
	}
	fmt.Println(tender)

	bids, err := s.bidRepo.GetBidsByTenderID(tenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bids: %w", err)
	}

	return bids, nil
}

func (s *BidService) DeleteBid(contractorID int, bidID int) error {
	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("Bid not found or access denied")
	}

	if bid.ContractorID != contractorID {
		return errors.New("Bid not found or access denied")
	}

	err = s.bidRepo.DeleteBid(bidID)
	if err != nil {
		return fmt.Errorf("failed to delete bid: %w", err)
	}

	return nil
}

func (s *BidService) UpdateBidStatus(contractorID int, bidID int, newStatus string) error {
	validStatuses := []string{model.BidStatusAwarded, model.BidStatusRejected, model.BidStatusRejected}
	if !contains(validStatuses, newStatus) {
		return errors.New("invalid status")
	}

	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", err)
	}

	if bid.ContractorID != contractorID {
		return errors.New("you do not have access to this bid")
	}

	bid.Status = newStatus
	err = s.bidRepo.UpdateBidStatus(bid)
	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (s *BidService) AwardBid(clientID int, tenderID int, bidID int) error {
	tender, err := s.tenderRepo.GetTenderByID(tenderID)
	if err != nil || tender.ClientID != clientID {
		return fmt.Errorf("Tender not found or access denied")
	}

	bid, err := s.bidRepo.GetBidByID(bidID)
	if err != nil {
		return fmt.Errorf("Bid not found")
	}
	fmt.Println(bid)
	err = s.bidRepo.AwardBid(bidID)
	if err != nil {
		return fmt.Errorf("failed to award bid: %w", err)
	}

	err = s.tenderRepo.UpdateTenderStatus(tenderID, "awarded")
	if err != nil {
		return fmt.Errorf("failed to update tender status: %w", err)
	}

	return nil
}
