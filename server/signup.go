package server

import (
	"net/http"

	"github.com/go-chi/httplog/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/joeychilson/links/db"
	"github.com/joeychilson/links/pages/signup"
	"github.com/joeychilson/links/pkg/validate"
)

func (s *Server) SignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signup.Page(signup.Props{FormProps: signup.FormProps{}}).Render(r.Context(), w)
	}
}

func (s *Server) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		oplog := httplog.LogEntry(ctx)

		email := r.FormValue("email")
		username := r.FormValue("username")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password-confirm")

		emailExists, err := s.queries.EmailExists(ctx, email)
		if err != nil {
			oplog.Error("error checking if email exists", "error", err)
			props := signup.Props{
				Error:     ErrorInternalServer,
				FormProps: signup.FormProps{},
			}
			s.RetargetPage(ctx, w, signup.Page(props))
			return
		}

		if emailExists {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validate.ValidationError{validate.EmailValue: validate.EmailExistsError},
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		validationError := validate.Email(email)
		if validationError != nil {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validationError,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		usernameExists, err := s.queries.UsernameExists(ctx, username)
		if err != nil {
			oplog.Error("error checking if username exists", "error", err)
			props := signup.Props{
				Error:     ErrorInternalServer,
				FormProps: signup.FormProps{},
			}
			s.RetargetPage(ctx, w, signup.Page(props))
			return
		}

		if usernameExists {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validate.ValidationError{validate.UsernameValue: validate.UsernameExistsError},
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		validationError = validate.Username(username)
		if validationError != nil {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validationError,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		validationError = validate.Password(password)
		if validationError != nil {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validationError,
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		if password != passwordConfirm {
			props := signup.FormProps{
				Email:    email,
				Username: username,
				Error:    validate.ValidationError{validate.PasswordValue: validate.PasswordMatchError},
			}
			signup.Form(props).Render(ctx, w)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			oplog.Error("error hashing password", "error", err)
			props := signup.Props{
				Error:     ErrorInternalServer,
				FormProps: signup.FormProps{},
			}
			s.RetargetPage(ctx, w, signup.Page(props))
			return
		}

		userID, err := s.queries.CreateUser(ctx, db.CreateUserParams{
			Email:    email,
			Username: username,
			Password: string(hashedPassword),
		})
		if err != nil {
			oplog.Error("error creating user", "error", err)
			props := signup.Props{
				Error:     ErrorInternalServer,
				FormProps: signup.FormProps{},
			}
			s.RetargetPage(ctx, w, signup.Page(props))
			return
		}

		err = s.sessionManager.Set(w, r, userID)
		if err != nil {
			oplog.Error("failed to set session", "error", err)
			props := signup.Props{
				Error:     ErrorInternalServer,
				FormProps: signup.FormProps{},
			}
			s.RetargetPage(ctx, w, signup.Page(props))
			return
		}

		oplog.Info("user signed up", "user_id", userID.String())
		s.Redirect(w, r, "/")
	}
}
