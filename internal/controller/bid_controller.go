package controller

import "tender-managment/internal/service"

var (
	bidService *service.BidService
)

func SetBidService(bidSer *service.BidService) {
	bidService = bidSer
}
