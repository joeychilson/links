package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/joeychilson/flixmetrics/database"
	"github.com/joeychilson/flixmetrics/server"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)
	server := server.New(queries)

	log.Println("Serving application @ http://localhost:8080")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
