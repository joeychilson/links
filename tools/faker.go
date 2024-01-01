package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/joeychilson/links/db"
)

type User struct {
	Email    string `fake:"{email}"`
	Username string `fake:"{username}"`
	Password string `fake:"{password}"`
}

type Link struct {
	UserID uuid.UUID `fake:"-"`
	Title  string    `fake:"{sentence:5}"`
	Url    string    `fake:"{url}"`
}

type Comment struct {
	UserID  uuid.UUID `fake:"-"`
	LinkID  uuid.UUID `fake:"-"`
	Content string    `fake:"{sentence:50}"`
}

type Reply struct {
	UserID  uuid.UUID `fake:"-"`
	Comment uuid.UUID `fake:"-"`
	Content string    `fake:"{sentence:50}"`
}

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	if len(os.Args) != 3 {
		log.Fatal("usage: faker <data_type> <count>")
	}

	dataType := os.Args[1]
	if dataType != "users" && dataType != "links" && dataType != "likes" && dataType != "comments" && dataType != "replies" && dataType != "votes" {
		log.Fatal("data must be one of: users, links, likes, comments, replies, votes")
	}

	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("count must be an integer")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	log.Println("Creating fake data...")
	switch dataType {
	case "users":
		createUsers(ctx, queries, count)
	case "links":
		createLinks(ctx, queries, count)
	case "likes":
		createLikes(ctx, queries, count)
	case "comments":
		createComments(ctx, queries, count)
	case "replies":
		createReplies(ctx, queries, count)
	case "votes":
		createVotes(ctx, queries, count)
	}
}

func createUsers(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating users...")
	for i := 0; i < count; i++ {
		user := User{}
		gofakeit.Struct(&user)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		_, err = queries.CreateUser(ctx, db.CreateUserParams{
			Email:    user.Email,
			Username: user.Username,
			Password: string(hashedPassword),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createLinks(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating links...")
	users, err := queries.UserList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		link := Link{}
		gofakeit.Struct(&link)

		user := users[gofakeit.Number(0, len(users)-1)]

		_, err = queries.CreateLink(ctx, db.CreateLinkParams{
			UserID: user.ID,
			Title:  link.Title,
			Url:    link.Url,
			Slug:   xid.New().String(),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createLikes(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating likes...")
	users, err := queries.UserList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	links, err := queries.LinkList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		link := links[gofakeit.Number(0, len(links)-1)]

		err = queries.CreateLike(ctx, db.CreateLikeParams{
			UserID: user.ID,
			LinkID: link.ID,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createComments(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating comments...")
	users, err := queries.UserList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	links, err := queries.LinkList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		comment := Comment{}
		gofakeit.Struct(&comment)

		user := users[gofakeit.Number(0, len(users)-1)]
		link := links[gofakeit.Number(0, len(links)-1)]

		err = queries.CreateComment(ctx, db.CreateCommentParams{
			UserID:  user.ID,
			LinkID:  link.ID,
			Content: comment.Content,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createReplies(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating replies...")
	users, err := queries.UserList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	comments, err := queries.CommentList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < count; i++ {
		reply := Reply{}
		gofakeit.Struct(&reply)

		user := users[gofakeit.Number(0, len(users)-1)]
		comment := comments[gofakeit.Number(0, len(comments)-1)]

		err = queries.CreateReply(ctx, db.CreateReplyParams{
			UserID:   user.ID,
			LinkID:   comment.LinkID,
			ParentID: comment.ID,
			Content:  reply.Content,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createVotes(ctx context.Context, queries *db.Queries, count int) {
	log.Println("Creating votes...")
	users, err := queries.UserList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	comments, err := queries.CommentList(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < count; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		comment := comments[gofakeit.Number(0, len(comments)-1)]
		vote := 0
		if gofakeit.Bool() {
			vote = 1
		} else {
			vote = -1
		}
		err = queries.CreateVote(ctx, db.CreateVoteParams{
			UserID:    user.ID,
			CommentID: comment.ID,
			Vote:      int32(vote),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
