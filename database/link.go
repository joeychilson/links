package database

import (
	"context"

	"github.com/google/uuid"
)

type CreateLinkParams struct {
	UserID uuid.UUID
	Title  string
	Url    string
}

func (q *Queries) CreateLink(ctx context.Context, arg CreateLinkParams) error {
	query := "INSERT INTO links (user_id, title, url) VALUES ($1, $2, $3)"
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.Title, arg.Url)
	return err
}

type LinkVoteParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
	Vote   int32
}

func (q *Queries) LinkVote(ctx context.Context, arg LinkVoteParams) error {
	query := `
        WITH existing_vote AS (
            SELECT vote FROM link_votes WHERE user_id = $1 AND link_id = $2 FOR UPDATE
        ), deleted AS (
            DELETE FROM link_votes WHERE user_id = $1 AND link_id = $2 AND vote = $3
        ), updated AS (
            UPDATE link_votes SET vote = $3 WHERE user_id = $1 AND link_id = $2 AND vote != $3
        )
        INSERT INTO link_votes (user_id, link_id, vote)
        SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM existing_vote);
    `
	_, err := q.db.Exec(ctx, query, arg.UserID, arg.LinkID, arg.Vote)
	return err
}

func (q *Queries) ScoreLink(ctx context.Context, linkID uuid.UUID) (int64, error) {
	query := "SELECT COALESCE(SUM(vote), 0) FROM link_votes WHERE link_id = $1"
	row := q.db.QueryRow(ctx, query, linkID)
	var score int64
	err := row.Scan(&score)
	return score, err
}
