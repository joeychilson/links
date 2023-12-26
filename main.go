package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/internal/session"
	"github.com/joeychilson/links/server"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")

	err := database.Migrate(databaseURL)
	if err != nil {
		log.Fatalf("Unable to migrate database: %v\n", err)
	}

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v\n", err)
	}
	defer dbpool.Close()

	queries := database.New(dbpool)

	encryptionKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SECURE_COOKIE_ENCRYPTION_KEY"))
	validationKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SECURE_COOKIE_VALIDATION_KEY"))

	cookie := securecookie.New(encryptionKey, validationKey)
	sessionManager := session.NewManager(cookie, queries)

	server := server.New(queries, sessionManager)

	log.Println("Serving lixy application @ http://localhost:8080")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
