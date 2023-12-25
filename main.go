package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joeychilson/inquire/database"
	"github.com/joeychilson/inquire/server"
)

func main() {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := database.New(dbpool)
	server := server.New(queries)

	log.Println("Serving inquire application @ http://localhost:8080")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
