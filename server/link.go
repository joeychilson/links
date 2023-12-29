package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/link"
)

func (s *Server) Link() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkID := r.URL.Query().Get("id")

		if linkID == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			oplog.Error("failed to parse link id", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		linkRow, err := s.queries.Link(ctx, database.LinkParams{
			UserID: userID,
			LinkID: linkUUID,
		})
		if err != nil {
			oplog.Error("failed to get link", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		commentRows, err := s.queries.CommentFeed(ctx, database.CommentFeedParams{
			UserID: userID,
			LinkID: linkUUID,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			oplog.Error("failed to get comment feed", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		oplog.Info("link page loaded", "link_id", linkID, "comments", len(commentRows), "vote", linkRow.UserVoted)
		link.Page(link.Props{User: user, Link: linkRow, Comments: commentRows}).Render(ctx, w)
	}
}

func (s *Server) LinkVote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkID := r.URL.Query().Get("link_id")
		voteDir := r.URL.Query().Get("vote")

		redirectURL := r.URL.Query().Get("redirect_url")
		if redirectURL == "" {
			redirectURL = "/"
		}

		if linkID == "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			oplog.Error("failed to parse link id", "error", err)
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		var vote int32
		if voteDir == "up" {
			vote = 1
		} else if voteDir == "down" {
			vote = -1
		} else {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		err = s.queries.LinkVote(ctx, database.LinkVoteParams{
			UserID: user.ID,
			LinkID: linkUUID,
			Vote:   vote,
		})
		if err != nil {
			oplog.Error("failed to vote on link", "error", err)
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		oplog.Info("user voted on link", "link_id", linkID)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
