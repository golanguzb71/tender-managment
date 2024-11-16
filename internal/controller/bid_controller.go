package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
)

var (
	bidService *service.BidService
)

// SetBidService sets the global bid service for handling requests.
func SetBidService(bidSer *service.BidService) {
	bidService = bidSer
}

// CreateBidHandler handles the creation of a new bid.
func CreateBidHandler(c *gin.Context) {
	var bid model.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the bid using the service layer
	bidID, err := bidService.CreateBid(bid.TenderID, bid.ContractorID, bid.Price, bid.DeliveryTime, bid.Comments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created bid ID
	c.JSON(http.StatusCreated, gin.H{"bid_id": bidID})
}

// GetBidsByTenderIDHandler retrieves all bids for a given tender.
func GetBidsByTenderIDHandler(c *gin.Context) {
	tenderID, err := strconv.Atoi(c.DefaultQuery("tender_id", "0"))
	if err != nil || tenderID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tender ID"})
		return
	}

	// Get bids using the service layer
	bids, err := bidService.GetBidsByTenderID(tenderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of bids
	c.JSON(http.StatusOK, bids)
}

// GetBidByIDHandler retrieves a single bid by its ID.
func GetBidByIDHandler(c *gin.Context) {
	bidID, err := strconv.Atoi(c.DefaultQuery("bid_id", "0"))
	if err != nil || bidID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bid ID"})
		return
	}

	// Get the bid using the service layer
	bid, err := bidService.GetBidByID(bidID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If no bid found, return a 404
	if bid == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "bid not found"})
		return
	}

	// Return the bid details
	c.JSON(http.StatusOK, bid)
}

// DeleteBidHandler deletes a bid by its ID.
func DeleteBidHandler(c *gin.Context) {
	bidID, err := strconv.Atoi(c.DefaultQuery("bid_id", "0"))
	if err != nil || bidID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bid ID"})
		return
	}

	// Delete the bid using the service layer
	if err := bidService.DeleteBid(bidID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UpdateBidStatusHandler updates the status of a bid.
func UpdateBidStatusHandler(c *gin.Context) {
	bidID, err := strconv.Atoi(c.DefaultQuery("bid_id", "0"))
	if err != nil || bidID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bid ID"})
		return
	}

	status := c.DefaultQuery("status", "")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	// Update the status of the bid using the service layer
	if err := bidService.UpdateBidStatus(bidID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
