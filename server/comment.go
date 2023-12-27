package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/joeychilson/links/database"
)

func (s *Server) Comment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		user := s.UserFromContext(r.Context())

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

		err = s.queries.CreateComment(r.Context(), database.CreateCommentParams{
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

func (s *Server) CommentVote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		user := s.UserFromContext(r.Context())

		commentID := r.URL.Query().Get("comment_id")
		voteDir := r.URL.Query().Get("vote")

		redirectURL := r.URL.Query().Get("redirect_url")
		if redirectURL == "" {
			redirectURL = "/"
		}

		if commentID == "" {
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		commentUUID, err := uuid.Parse(commentID)
		if err != nil {
			oplog.Error("failed to parse comment id", "error", err)
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

		err = s.queries.CommentVote(r.Context(), database.CommentVoteParams{
			UserID:    user.ID,
			CommentID: commentUUID,
			Vote:      vote,
		})
		if err != nil {
			oplog.Error("failed to vote on comment", "error", err)
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}

		oplog.Info("user voted on comment", "comment_id", commentUUID)
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
