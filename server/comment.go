package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/joeychilson/links/components/reply"
	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/link"
)

func (s *Server) Comment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkID := r.FormValue("link_id")
		parentID := r.FormValue("parent_id")
		content := r.FormValue("content")

		if linkID == "" {
			s.Redirect(w, r, "/")
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			oplog.Error("failed to parse link id", "error", err)
			s.Redirect(w, r, "/")
			return
		}

		if content == "" {
			s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
			return
		}

		var userID uuid.UUID
		if user != nil {
			userID = user.ID
		} else {
			s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
			return
		}

		var parentUUID uuid.NullUUID
		if parentID != "" {
			parentUUID.UUID, err = uuid.Parse(parentID)
			if err != nil {
				oplog.Error("failed to parse parent id", "error", err)
				s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
				return
			}
			parentUUID.Valid = true
		} else {
			parentUUID.Valid = false
		}

		err = s.queries.CreateComment(ctx, database.CreateCommentParams{
			UserID:   userID,
			LinkID:   linkUUID,
			ParentID: parentUUID,
			Content:  content,
		})
		if err != nil {
			oplog.Error("failed to create comment", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
			return
		}

		commentFeed, err := s.queries.CommentFeed(ctx, database.CommentFeedParams{
			LinkID: linkUUID,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			oplog.Error("failed to get comment feed", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
			return
		}

		oplog.Info("user created comment", "link_id", linkID)
		w.Header().Set("HX-Refresh", "true")
		link.CommentFeed(link.CommentFeedProps{User: user, LinkID: linkID, CommentFeed: commentFeed}).Render(ctx, w)
	}
}

func (s *Server) CommentReply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkID := r.URL.Query().Get("link_id")
		if linkID == "" {
			s.Redirect(w, r, "/")
			return
		}

		commentID := r.URL.Query().Get("comment_id")
		if commentID == "" {
			s.Redirect(w, r, fmt.Sprintf("/link?id=%s", linkID))
			return
		}

		reply.Component(reply.Props{LinkID: linkID, CommentID: commentID}).Render(r.Context(), w)
	}
}
