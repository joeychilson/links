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
	"github.com/joeychilson/links/pages/link"
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
			props := &link.CommentTextboxProps{LinkSlug: slug, Error: "Please enter a comment"}
			link.CommentTextbox(props).Render(ctx, w)
			return
		}

		commentID, err := s.queries.CreateComment(ctx, db.CreateCommentParams{
			UserID:  user.ID,
			LinkID:  linkID,
			Content: content,
		})
		if err != nil {
			oplog.Error("error creating comment", err)
			props := &link.CommentTextboxProps{LinkSlug: slug, Content: content, Error: "Sorry, something went wrong"}
			link.CommentTextbox(props).Render(ctx, w)
			return
		}

		commentRow, err := s.queries.Comment(ctx, db.CommentParams{
			ID:     commentID,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comment", err)
			props := &link.CommentTextboxProps{LinkSlug: slug, Content: content, Error: "Sorry, something went wrong"}
			link.CommentTextbox(props).Render(ctx, w)
			return
		}

		oplog.Info("comment created", "slug", slug)
		link.CommentTextbox(&link.CommentTextboxProps{LinkSlug: slug, Content: ""}).Render(ctx, w)
		comment.Component(&comment.Props{User: user, CommentRow: commentRow}).Render(ctx, w)
	}
}

func (s *Server) ReplyTextbox() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		linkSlug := chi.URLParam(r, "slug")
		commentID := chi.URLParam(r, "commentID")
		reply.Component(&reply.Props{LinkSlug: linkSlug, CommentID: commentID}).Render(r.Context(), w)
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
			props := &reply.Props{LinkSlug: linkSlug, CommentID: commentID, Error: "Please enter a reply"}
			reply.Component(props).Render(ctx, w)
			return
		}

		replyID, err := s.queries.CreateReply(ctx, db.CreateReplyParams{
			UserID:   user.ID,
			LinkID:   linkID,
			ParentID: commentUUID,
			Content:  content,
		})
		if err != nil {
			oplog.Error("error creating reply", err)
			props := &reply.Props{LinkSlug: linkSlug, CommentID: commentID, Error: "Sorry, something went wrong"}
			reply.Component(props).Render(ctx, w)
			return
		}

		commentRow, err := s.queries.Comment(ctx, db.CommentParams{
			ID:     replyID,
			UserID: user.ID,
		})
		if err != nil {
			oplog.Error("error getting comment", err)
			props := &reply.Props{LinkSlug: linkSlug, CommentID: commentID, Error: "Sorry, something went wrong"}
			reply.Component(props).Render(ctx, w)
			return
		}

		oplog.Info("reply created", "slug", linkSlug)
		comment.Component(&comment.Props{User: user, CommentRow: commentRow}).Render(ctx, w)
	}
}
