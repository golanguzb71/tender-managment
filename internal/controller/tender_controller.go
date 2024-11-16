package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
	"time"
)

var (
	tenderService *service.TenderService
)

func SetTenderService(tenderSer *service.TenderService) {
	tenderService = tenderSer

}

func CreateTender(c *gin.Context) {
	var payload model.CreateTender

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	if payload.Budget < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender data"})
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender data"})
		return
	}

	tender := model.Tender{
		Title:       payload.Title,
		Description: payload.Description,
		Deadline:    parsedTime,
		Budget:      payload.Budget,
		Attachment:  payload.Attachment,
		ClientID:    c.GetInt("user_id"),
		Status:      "open",
	}

	createdTender, err := tenderService.CreateTender(&tender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create tender", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTender)
}

func ListTenders(c *gin.Context) {
	clientID := c.GetInt("user_id")
	tenders, err := tenderService.ListTenders(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tenders"})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

func UpdateTenderStatus(c *gin.Context) {
	tenderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender ID"})
		return
	}

	var statusUpdate model.UpdateTenderStatusRequest
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	clientID := c.GetInt("user_id")
	if err := tenderService.UpdateTenderStatus(clientID, tenderID, statusUpdate.Status); err != nil {
		switch err.Error() {
		case "tender not found":
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found"})
		case "invalid status":
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender status"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update tender"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tender status updated"})
}

func DeleteTender(c *gin.Context) {
	tenderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender ID"})
		return
	}

	clientID := c.GetInt("user_id")
	if err := tenderService.DeleteTender(clientID, tenderID); err != nil {
		if err.Error() == "tender not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete tender"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tender deleted successfully"})
}
