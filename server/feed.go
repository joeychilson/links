package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/joeychilson/links/database"
	"github.com/joeychilson/links/templates/pages/feed"
	"github.com/joeychilson/links/templates/pages/new"
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

		linkFeedRows, err := s.queries.LinkFeed(r.Context(), database.LinkFeedParams{
			UserID: user.ID,
			Limit:  25,
			Offset: int32((page - 1) * 25),
		})
		if err != nil {
			log.Println(err)
			feed.Page(feed.Props{User: user}).Render(r.Context(), w)
			return
		}

		for i, linkFeedRow := range linkFeedRows {
			fmt.Println(i, linkFeedRow.UserVoted)
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

func (s *Server) Vote() http.HandlerFunc {
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

		userVoted, err := s.queries.UserVoted(r.Context(), database.UserVotedParams{
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
			err = s.queries.CreateVote(r.Context(), database.CreateVoteParams{
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
