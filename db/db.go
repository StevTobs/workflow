package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {

	DATABASE_URL := os.Getenv("DATABASE_URL")

	var err error
	DB, err = pgxpool.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	fmt.Println("Connected to the database!")
}
