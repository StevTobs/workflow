package models

import (
	// "context"
	// "fmt"
	// "item-workflow-system/db"
	// "log"
	// "github.com/jackc/pgx/v4/pgxpool"
)

// var DB *pgxpool.Pool// "github.com/lib/pq"


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
// 	err := DB.QueryRow(context.Background(), query, item.Title, item.Amount, item.Quantity, item.Status, item.Owner_id).Scan(&newId)
// 	if err != nil {
// 		log.Printf("Failed to insert new item: %v", err)
// 		return 0, err
// 	}

//		fmt.Println("New item created with ID:", newId)
//		return newId, nil
//	}

// VerifyCredentials checks if the username and password are valid
// func VerifyCredentials(username, password string) (bool, error) {
// 	// Implement your user credential verification logic here.
// 	// For example, you might query your database to check if the username and password are correct.
// 	// This is just a placeholder implementation.

// 	// Example: check if username and password match hardcoded values
// 	if username == "admin" && password == "password" {
// 		return true, nil
// 	}

// 	return false, nil
// }
