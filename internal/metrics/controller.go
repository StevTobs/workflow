package metrics

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)
type Message struct {
	Text string `json:"text"`
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow CORS
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		response := Message{Text: "Hello from Go server!"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ReadAllItems_(c *gin.Context) {
	items, err := GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Dummy function to simulate getting an item by ID (replace with your actual implementation)
var ErrItemNotFound = errors.New("item not found")
func ReadItemByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	item, err := GetItemByID(id)
	if err != nil {
		if err == ErrItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		}
		return
	}

	c.JSON(http.StatusOK, item)
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