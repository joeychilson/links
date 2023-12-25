package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joeychilson/lixy/database"
	"github.com/joeychilson/lixy/server"
)

func main() {
	ctx := context.Background()

	err := database.Migrate(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to migrate database: %v\n", err)
	}

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v\n", err)
	}
	defer dbpool.Close()

	queries := database.New(dbpool)
	server := server.New(queries)

	log.Println("Serving lixy application @ http://localhost:8080")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
