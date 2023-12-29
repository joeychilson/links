package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/link"
)

func (s *Server) Comment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkID := r.FormValue("link_id")
		content := r.FormValue("content")

		if linkID == "" {
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			oplog.Error("failed to parse link id", "error", err)
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
			return
		}

		if content == "" {
			w.Header().Set("HX-Redirect", fmt.Sprintf("/link?id=%s", linkID))
			w.WriteHeader(http.StatusOK)
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			userID = uuid.Nil
		}

		err = s.queries.CreateComment(ctx, database.CreateCommentParams{
			UserID:  userID,
			LinkID:  linkUUID,
			Content: content,
		})
		if err != nil {
			oplog.Error("failed to create comment", "error", err)
			w.Header().Set("HX-Redirect", fmt.Sprintf("/link?id=%s", linkID))
			w.WriteHeader(http.StatusOK)
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
			w.Header().Set("HX-Redirect", fmt.Sprintf("/link?id=%s", linkID))
			w.WriteHeader(http.StatusOK)
			return
		}

		oplog.Info("user created comment", "link_id", linkID)
		link.CommentFeed(link.CommentFeedProps{User: user, LinkID: linkID, Comments: commentRows}).Render(ctx, w)
	}
}
