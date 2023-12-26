package sessions

import (
	"encoding/base64"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/joeychilson/lixy/database"
)

const CookieName = "session"

type Manager struct {
	cookie  *securecookie.SecureCookie
	queries *database.Queries
}

func New(cookie *securecookie.SecureCookie, queries *database.Queries) *Manager {
	return &Manager{
		cookie:  cookie,
		queries: queries,
	}
}

func (m *Manager) Set(w http.ResponseWriter, r *http.Request, userID int32) error {
	token, err := m.queries.CreateUserToken(r.Context(), database.CreateUserTokenParams{
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

func (m *Manager) Delete(w http.ResponseWriter, r *http.Request) error {
	cookie, err := m.Get(r)
	if err != nil {
		return err
	}
	err = m.queries.DeleteUserToken(r.Context(), database.DeleteUserTokenParams{
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
