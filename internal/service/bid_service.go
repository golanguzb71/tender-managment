package service

import repository "tender-managment/internal/db/repo"

type BidService struct {
	repo *repository.BidRepository
}

func NewBidService(repo *repository.BidRepository) *BidService {
	return &BidService{repo: repo}
}
