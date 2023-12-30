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

type LinkScoreAndUserVoteParams struct {
	UserID uuid.UUID
	LinkID uuid.UUID
}

type LinkScoreAndUserVoteRow struct {
	Score    int64
	UserVote int32
}

func (q *Queries) LinkScoreAndUserVote(ctx context.Context, args LinkScoreAndUserVoteParams) (LinkScoreAndUserVoteRow, error) {
	query := `
		SELECT 
			COALESCE(SUM(vote), 0) AS score,
			COALESCE((SELECT vote FROM link_votes WHERE user_id = $1 AND link_id = $2), 0) AS user_vote
		FROM 
			link_votes
		WHERE 
			link_id = $2
	`
	var row LinkScoreAndUserVoteRow
	if err := q.db.QueryRow(ctx, query, args.UserID, args.LinkID).Scan(&row.Score, &row.UserVote); err != nil {
		return LinkScoreAndUserVoteRow{}, err
	}
	return row, nil
}
