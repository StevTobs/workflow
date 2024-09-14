package models

import (
	"context"
	"fmt"
	"item-workflow-system/db"
	"log"
	// "github.com/lib/pq"
)

type Item struct {
	Id       int    `json:"id"`
	Title    string `json:"title" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
	Status   string `json:"status"`
	Owner_id int    `json:"owner_id"`
}

// 1. ระบบสามารถเพมิ ขอ้ มลู การเบกิ งบใหมไ่ ด้
// Create a POST Endpoint
// func CreateItem(item Item) (int, error) {
// 	var newId int
// 	query := `
// 		INSERT INTO items (title, amount, quantity, status, owner_id)
// 		VALUES ($1, $2, $3, $4, $5)
// 		RETURNING id
// 	`
// 	err := db.DB.QueryRow(context.Background(), query, item.Title, item.Amount, item.Quantity, item.Status, item.Owner_id).Scan(&newId)
// 	if err != nil {
// 		log.Printf("Failed to insert new item: %v", err)
// 		return 0, err
// 	}

//		fmt.Println("New item created with ID:", newId)
//		return newId, nil
//	}

// VerifyCredentials checks if the username and password are valid
func VerifyCredentials(username, password string) (bool, error) {
	// Implement your user credential verification logic here.
	// For example, you might query your database to check if the username and password are correct.
	// This is just a placeholder implementation.

	// Example: check if username and password match hardcoded values
	if username == "admin" && password == "password" {
		return true, nil
	}

	return false, nil
}


func CreateItem(title string, amount int, quantity int) (int, error) {
	// Assuming you have a global DB variable initialized elsewhere
	var id int

	// Define the query to insert the new item
	query := `INSERT INTO items (title, amount, quantity, status, owner_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	// Execute the query
	err := db.DB.QueryRow(context.Background(), query, title, amount, quantity, "PENDING", 1).Scan(&id)
	if err != nil {
		log.Println("Error inserting item:", err)
		return 0, err
	}

	return id, nil
}

// 2.ระบบสามารถดขู อ้ มลู การเบกิ งบทงั หมดได้
// Set Up a GET /items Endpoint
func GetAllItems() ([]Item, error) {

	rows, err := db.DB.Query(context.Background(), "SELECT id, title, amount, quantity, status, owner_id FROM items")

	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
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
func GetItemByID(id int) (Item, error) {
	var item Item
	query := `
		SELECT id, title, amount, quantity, status, owner_id 
		FROM items 
		WHERE id = $1
	`
	err := db.DB.QueryRow(context.Background(), query, id).Scan(&item.Id, &item.Title, &item.Amount, &item.Quantity, &item.Status, &item.Owner_id)
	if err != nil {
		log.Printf("Failed to get item with ID %d: %v", id, err)
		return item, err
	}

	fmt.Printf("Item with ID %d retrieved successfully!\n", id)
	return item, nil
}

// 4. ระบบสามารถปรบั เปลียน/แก้ไขข้อมูล การเบกิ งบได้
// UpdateItem updates the item's details by its ID
func UpdateItem(id int, title string, amount int, quantity int) error {
	// Define the query to update the item
	query := `UPDATE items SET title = $1, amount = $2, quantity = $3 WHERE id = $4`

	// Execute the query
	_, err := db.DB.Exec(context.Background(), query, title, amount, quantity, id)
	if err != nil {
		log.Println("Error updating item:", err)
		return err
	}

	return nil

}

// 5. ระบบสามารถปรับเปลี่ยนแก้ไขข้อมูลสถานะการเบิกงบได้
// UpdateItemStatus updates the status of an item by its ID
func UpdateItemStatus(id int, newStatus string) error {
	// Define the query to update the item status
	query := `UPDATE items SET status = $1 WHERE id = $2`

	// Execute the query
	_, err := db.DB.Exec(context.Background(), query, newStatus, id)
	if err != nil {
		log.Println("Error updating item status:", err)
		return err
	}

	return nil
}

// 6. ระบบสามารถลบขอ้ มลู การเบกิ งบได้
// DeleteItem deletes an item from the database by its ID
func DeleteItem(id int) error {
	// Define the query to delete the item
	query := `DELETE FROM items WHERE id = $1`

	// Execute the query
	_, err := db.DB.Exec(context.Background(), query, id)
	if err != nil {
		log.Println("Error deleting item:", err)
		return err
	}

	return nil
}
