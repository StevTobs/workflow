package metrics

import (
	"context"
	"fmt"
	model "item-workflow-system/internal/model"
	"log"
	"net/http"
	"os"
	"strconv"

	// "github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	// Load environment variables from .env file (if you have one)
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("Error loading .env file")
	// }

	DATABASE_URL := os.Getenv("DATABASE_URL")
	PORT := os.Getenv("PORT")


	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if PORT == "" {
		log.Fatal("PORT is not set")
	}

	var err error
	DB, err = pgxpool.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	fmt.Printf("Connected to the database on port %s!\n", PORT)
}


func CreateItem(title string, amount int, quantity int) (int, error) {
	// Assuming you have a global DB variable initialized elsewhere
	var id int

	// Define the query to insert the new item
	query := `INSERT INTO items (title, amount, quantity, status, owner_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Execute the query
	err := DB.QueryRow(context.Background(), query, title, amount, quantity, "PENDING", 1).Scan(&id)
	if err != nil {
		log.Println("Error inserting item:", err)
		return 0, err
	}

	return id, nil
}

// 2.ระบบสามารถดขู อ้ มลู การเบกิ งบทงั หมดได้
// Set Up a GET /items Endpoint


func GetAllItems() ([]model.Item, error) {

	rows, err := DB.Query(context.Background(), "SELECT id, title, amount, quantity, status, owner_id FROM items")

	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []model.Item

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Title, &item.Amount, &item.Quantity, &item.Status, &item.Owner_id)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}
		items = append(items, item)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return items, nil

}

// 3.ระบบสามารถดขู อ้ มลู การเบกิ งบทลี ะรายการได้
func GetItemByID(id int) (model.Item, error) {
	var item model.Item
	query := `
		SELECT id, title, amount, quantity, status, owner_id 
		FROM items 
		WHERE id = $1
	`
	err := DB.QueryRow(context.Background(), query, id).Scan(&item.Id, &item.Title, &item.Amount, &item.Quantity, &item.Status, &item.Owner_id)
	if err != nil {
		log.Printf("Failed to get item with ID %d: %v", id, err)
		return item, err
	}

	fmt.Printf("Item with ID %d retrieved successfully!\n", id)
	return item, nil
}

// 4. ระบบสามารถปรบั เปลียน/แก้ไขข้อมูล การเบกิ งบได้
// UpdateItem updates the item's details by its ID
// func UpdateItem(id int, title string, amount int, quantity int) error {
// 	// Define the query to update the item
// 	query := `UPDATE items SET title = $1, amount = $2, quantity = $3 WHERE id = $4`

// 	// Execute the query
// 	_, err := DB.Exec(context.Background(), query, title, amount, quantity, id)
// 	if err != nil {
// 		log.Println("Error updating item:", err)
// 		return err
// 	}

// 	return nil

// }
func UpdateItem(id int, title string, amount int, quantity int) error {
	// Define the query to update the item
	query := `UPDATE items SET title = $1, amount = $2, quantity = $3 WHERE id = $4`

	// Execute the query with the provided parameters
	_, err := DB.Exec(context.Background(), query, title, amount, quantity, id)
	if err != nil {
		log.Printf("Error updating item with ID %d: %v", id, err)
		return fmt.Errorf("failed to update item with ID %d: %w", id, err) // Wrap error for better context
	}

	return nil
}


// 5. ระบบสามารถปรับเปลี่ยนแก้ไขข้อมูลสถานะการเบิกงบได้
// UpdateItemStatus updates the status of an item by its ID
func UpdateItemStatus(id int, newStatus string) error {
	// Define the query to update the item status
	query := `UPDATE items SET status = $1 WHERE id = $2`

	// Execute the query
	_, err := DB.Exec(context.Background(), query, newStatus, id)
	if err != nil {
		log.Println("Error updating item status:", err)
		return err
	}

	return nil
}
func PatchItemStatusHandler(c *gin.Context) {
	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	// Bind JSON request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Extract the ID from the URL parameters
	// id := c.Param("id")
	idStr := c.Param("id") // Example of how you might get the ID from the URL
	id, err := strconv.Atoi(idStr)

	// Call the function to update the item's status
	err = UpdateItemStatus(id, requestBody.Status)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item status updated successfully"})
}


// 6. ระบบสามารถลบขอ้ มลู การเบกิ งบได้
// DeleteItem deletes an item from the database by its ID
func DeleteItem(id int) error {
	// Define the query to delete the item
	query := `DELETE FROM items WHERE id = $1`

	// Execute the query
	_, err := DB.Exec(context.Background(), query, id)
	if err != nil {
		log.Println("Error deleting item:", err)
		return err
	}

	return nil
}
