package users

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joeychilson/lixy/pkg/context"
)

const ContextKey context.ContextKey = "user"

type User struct {
	ID       pgtype.UUID
	Email    string
	Username string
}
