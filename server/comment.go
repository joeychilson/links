package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"

	"github.com/joeychilson/links/components/comment"
	"github.com/joeychilson/links/components/reply"
	"github.com/joeychilson/links/db"
)

func (s *Server) Comment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")
		content := r.FormValue("content")

		linkID, err := s.queries.LinkIDBySlug(ctx, slug)
		if err != nil {
			oplog.Error("error getting link id", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		if content == "" {
			props := comment.TextboxProps{LinkSlug: slug, Error: "Please enter a comment"}
			comment.Textbox(props).Render(ctx, w)
			return
		}

		err = s.queries.CreateComment(ctx, db.CreateCommentParams{
			UserID:  user.ID,
			LinkID:  linkID,
			Content: content,
		})
		if err != nil {
			oplog.Error("error creating comment", err)
			props := comment.TextboxProps{LinkSlug: slug, Content: content, Error: "Sorry, something went wrong"}
			comment.Textbox(props).Render(ctx, w)
			return
		}

		oplog.Info("comment created", "slug", slug)
		s.RefreshPage(w, r)
	}
}

func (s *Server) ReplyTextbox() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkSlug := chi.URLParam(r, "slug")
		commentID := chi.URLParam(r, "commentID")
		reply.Component(reply.Props{LinkSlug: linkSlug, CommentID: commentID}).Render(r.Context(), w)
	}
}

func (s *Server) Reply() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		linkSlug := chi.URLParam(r, "slug")
		commentID := chi.URLParam(r, "commentID")
		content := r.FormValue("content")

		linkID, err := s.queries.LinkIDBySlug(ctx, linkSlug)
		if err != nil {
			oplog.Error("error getting link id", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", linkSlug))
			return
		}

		commentUUID, err := uuid.Parse(commentID)
		if err != nil {
			oplog.Error("failed to parse parent id", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", linkSlug))
			return
		}

		if content == "" {
			props := reply.Props{LinkSlug: linkSlug, CommentID: commentID, Error: "Please enter a reply"}
			reply.Component(props).Render(ctx, w)
			return
		}

		err = s.queries.CreateReply(ctx, db.CreateReplyParams{
			UserID:   user.ID,
			LinkID:   linkID,
			ParentID: commentUUID,
			Content:  content,
		})
		if err != nil {
			oplog.Error("error creating reply", err)
			props := reply.Props{LinkSlug: linkSlug, CommentID: commentID, Error: "Sorry, something went wrong"}
			reply.Component(props).Render(ctx, w)
			return
		}

		oplog.Info("reply created", "slug", linkSlug)
		s.RefreshPage(w, r)
	}
}

func (s *Server) Upvote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")
		commentID := chi.URLParam(r, "commentID")

		commentUUID, err := uuid.Parse(commentID)
		if err != nil {
			oplog.Error("failed to parse comment id", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		err = s.queries.CreateVote(ctx, db.CreateVoteParams{
			UserID:    user.ID,
			CommentID: commentUUID,
			Vote:      1,
		})
		if err != nil {
			oplog.Error("failed to vote on comment", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		commentRow, err := s.queries.Comment(ctx, db.CommentParams{
			ID:     commentUUID,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("failed to get comment", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		oplog.Info("upvoted", "comment_id", commentID)
		comment.Component(comment.Props{User: user, CommentRow: commentRow}).Render(ctx, w)
	}
}

func (s *Server) Downvote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)
		user := s.UserFromContext(ctx)
		slug := chi.URLParam(r, "slug")
		commentID := chi.URLParam(r, "commentID")

		commentUUID, err := uuid.Parse(commentID)
		if err != nil {
			oplog.Error("failed to parse comment id", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		err = s.queries.CreateVote(ctx, db.CreateVoteParams{
			UserID:    user.ID,
			CommentID: commentUUID,
			Vote:      -1,
		})
		if err != nil {
			oplog.Error("failed to vote on comment", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		commentRow, err := s.queries.Comment(ctx, db.CommentParams{
			ID:     commentUUID,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("failed to get comment", "error", err)
			s.Redirect(w, r, fmt.Sprintf("/%s", slug))
			return
		}

		oplog.Info("downvoted", "comment_id", commentID)
		comment.Component(comment.Props{User: user, CommentRow: commentRow}).Render(ctx, w)
	}
}
