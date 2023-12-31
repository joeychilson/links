package session

import (
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"

	"github.com/joeychilson/links/db"
)

// User represents a user in a session.
type User struct {
	ID       uuid.UUID
	Avatar   string
	Email    string
	Username string
}

// ContextKey is the key used to store the session in the request context.
type ContextKey string

const (
	CookieName            = "session"
	SessionKey ContextKey = "session"
)

// Manager represents a session manager.
type Manager struct {
	cookie  *securecookie.SecureCookie
	queries *db.Queries
}

// NewManager returns a new session manager.
func NewManager(cookie *securecookie.SecureCookie, queries *db.Queries) *Manager {
	return &Manager{
		cookie:  cookie,
		queries: queries,
	}
}

// Set sets a session for the given user.
func (m *Manager) Set(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	token, err := m.queries.CreateUserToken(r.Context(), db.CreateUserTokenParams{
		UserID:  userID,
		Token:   base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)),
		Context: CookieName,
	})
	if err != nil {
		return err
	}
	encoded, err := m.cookie.Encode(CookieName, token)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    encoded,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}

// Get gets the session for the given request.
func (m *Manager) Get(r *http.Request) (string, error) {
	var value string
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	err = m.cookie.Decode(CookieName, cookie.Value, &value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetUser gets the user for the given request.
func (m *Manager) GetUser(r *http.Request) (User, error) {
	ctx := r.Context()
	cookie, err := m.Get(r)
	if err != nil {
		return User{ID: uuid.Nil}, err
	}
	userID, err := m.queries.UserIDByToken(ctx, db.UserIDByTokenParams{
		Token:   cookie,
		Context: CookieName,
	})
	if err != nil {
		return User{ID: uuid.Nil}, err
	}
	userRow, err := m.queries.UserByID(ctx, userID)
	if err != nil {
		return User{ID: uuid.Nil}, err
	}
	return User{
		ID:       userRow.ID,
		Email:    userRow.Email,
		Username: userRow.Username,
	}, nil
}

// Delete deletes the session for the given request.
func (m *Manager) Delete(w http.ResponseWriter, r *http.Request) error {
	cookie, err := m.Get(r)
	if err != nil {
		return err
	}
	err = m.queries.DeleteUserToken(r.Context(), db.DeleteUserTokenParams{
		Token:   cookie,
		Context: CookieName,
	})
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  "",
		MaxAge: -1,
	})
	return nil
}
