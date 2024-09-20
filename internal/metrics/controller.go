package metrics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReadAllItems fetches all items from the model and returns them in a JSON response.
func ReadAllItems_(c *gin.Context) {
	items, err := GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}


func CreateItem_(c *gin.Context) {
	var requestBody struct {
		Title    string `json:"title" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newId, err := CreateItem(requestBody.Title, requestBody.Amount, requestBody.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item created successfully", "id": newId})
}

// Define the PUT endpoint to update an item by ID (protected)
func UpdateItem_(c *gin.Context) {
	var requestBody struct {
		Title    string `json:"title" binding:"required"`
		Amount   int    `json:"amount" binding:"required"`
		Quantity int    `json:"quantity" binding:"required"`
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = UpdateItem(id, requestBody.Title, requestBody.Amount, requestBody.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// Define the PATCH endpoint to update the status of an item by ID (protected)
func PatialUpdateItem_(c *gin.Context) {
	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Status != "PENDING" && requestBody.Status != "APPROVED" && requestBody.Status != "REJECTED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	err = UpdateItemStatus(id, requestBody.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item status updated successfully"})
}

// Define the DELETE endpoint to delete an item by ID (protected)
func RemoveItem_(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = DeleteItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	c.Status(http.StatusNoContent)
}