package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/joeychilson/links/database"
)

func (s *Server) Comment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)

		linkID := r.FormValue("link_id")
		content := r.FormValue("content")

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

		if content == "" {
			http.Redirect(w, r, "/link?id="+linkID, http.StatusFound)
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
			http.Redirect(w, r, "/link?id="+linkID, http.StatusFound)
			return
		}

		oplog.Info("user created comment", "link_id", linkID)
		http.Redirect(w, r, "/link?id="+linkID, http.StatusFound)
	}
}
