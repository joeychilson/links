package main

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/gorilla/securecookie"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pkg/session"
	"github.com/joeychilson/links/server"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	logger := httplog.NewLogger("links", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		JSON:             false,
		RequestHeaders:   false,
		ResponseHeaders:  false,
		MessageFieldName: "message",
		QuietDownPeriod:  10 * time.Second,
	})

	err := db.Migrate(databaseURL)
	if err != nil {
		slog.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	encryptionKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SECURE_COOKIE_ENCRYPTION_KEY"))
	validationKey, _ := base64.StdEncoding.DecodeString(os.Getenv("SECURE_COOKIE_VALIDATION_KEY"))

	cookie := securecookie.New(encryptionKey, validationKey)
	sessionManager := session.NewManager(cookie, queries)

	server := server.New(logger, queries, sessionManager)

	slog.Info("Starting links application @ http://localhost:8080")
	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
