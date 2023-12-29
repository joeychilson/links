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

		commentFeed, err := s.queries.CommentFeed(ctx, database.CommentFeedParams{
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

		oplog.Info("link page loaded", "link_id", linkID, "comments", len(commentFeed), "vote", linkRow.UserVoted)
		link.Page(link.Props{User: user, Link: linkRow, CommentFeed: commentFeed}).Render(ctx, w)
	}
}
