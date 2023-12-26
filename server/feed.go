package server

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/pages/feed"
	"github.com/joeychilson/links/pages/new"
)

func (s *Server) FeedPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())

		pageStr := r.URL.Query().Get("page")
		if pageStr == "" {
			pageStr = "1"
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		var linkFeedRows []database.LinkFeedRow

		if user != nil {
			linkFeedRows, err = s.queries.UserLinkFeed(r.Context(), database.UserLinkFeedParams{
				UserID: user.ID,
				Limit:  25,
				Offset: int32((page - 1) * 25),
			})
			if err != nil {
				feed.Page(feed.Props{User: user}).Render(r.Context(), w)
				return
			}
		} else {
			linkFeedRows, err = s.queries.LinkFeed(r.Context(), database.LinkFeedParams{
				Limit:  25,
				Offset: int32((page - 1) * 25),
			})
			if err != nil {
				feed.Page(feed.Props{User: user}).Render(r.Context(), w)
				return
			}
		}

		feed.Page(feed.Props{User: user, LinkFeedRows: linkFeedRows}).Render(r.Context(), w)
	}
}

func (s *Server) NewPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		new.Page(new.PageProps{User: user}).Render(r.Context(), w)
	}
}

func (s *Server) New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())

		title := r.FormValue("title")
		url := r.FormValue("url")

		if title == "" || url == "" {
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		err := s.queries.CreateLink(r.Context(), database.CreateLinkParams{
			UserID: user.ID,
			Title:  title,
			Url:    url,
		})
		if err != nil {
			new.Page(new.PageProps{User: user}).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Server) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := s.UserFromContext(r.Context())
		linkID := r.URL.Query().Get("link_id")

		if linkID == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		linkUUID, err := uuid.Parse(linkID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		userVoted, err := s.queries.UserLiked(r.Context(), database.UserLikedParams{
			UserID: user.ID,
			LinkID: linkUUID,
		})
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if userVoted {
			err = s.queries.DeleteVote(r.Context(), database.DeleteVoteParams{
				UserID: user.ID,
				LinkID: linkUUID,
			})
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		} else {
			err = s.queries.CreateLike(r.Context(), database.CreateLikeParams{
				UserID: user.ID,
				LinkID: linkUUID,
			})
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
